// +build sbs

package ads_b

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/jphastings/jan-poka/pkg/locator/common"
	"github.com/jphastings/jan-poka/pkg/math"
	"github.com/twuillemin/modes/pkg/bds/adsb"
	"github.com/twuillemin/modes/pkg/bds/bds05/fields"
	"github.com/twuillemin/modes/pkg/bds/bds05/messages"
	messages2 "github.com/twuillemin/modes/pkg/bds/bds08/messages"
	messages3 "github.com/twuillemin/modes/pkg/bds/bds09/messages"
	fields2 "github.com/twuillemin/modes/pkg/bds/bds61/fields"
	messages6 "github.com/twuillemin/modes/pkg/bds/bds61/messages"
	messages5 "github.com/twuillemin/modes/pkg/bds/bds62/messages"
	messages4 "github.com/twuillemin/modes/pkg/bds/bds65/messages"
	adsbReader "github.com/twuillemin/modes/pkg/bds/reader"
	"github.com/twuillemin/modes/pkg/geo"
	modeSMessages "github.com/twuillemin/modes/pkg/modes/messages"
	modeSReader "github.com/twuillemin/modes/pkg/modes/reader"
	"log"
	"net"
	"strings"
)

var _ common.LocationProvider = (*sbsLocationProvider)(nil)

var positionCache *sbsLocationProvider

type sbsLocationProvider struct {
	dataCache    map[uint32]planeData
	focusAddress uint32
	name         string
}

func init() {
	positionCache = &sbsLocationProvider{
		dataCache: map[uint32]planeData{},
	}
	if err := positionCache.connectDecoder("127.0.0.1:30002"); err != nil {
		log.Println("❌ Provider: ADS-B unavailable, no server running at localhost 30002")
		return
	}

	common.Providers[TYPE] = func() common.LocationProvider { return positionCache }
	log.Println("✅ Provider: ADS-B airplane positions available.")
}

func (lp *sbsLocationProvider) SetParams(decodeInto func(interface{}) error) error {
	loc := &params{}
	err := decodeInto(loc)
	if err != nil {
		return err
	}
	var address uint32
	var trackedFlights []string
	for addr, data := range lp.dataCache {
		if data.ident == "" {
			continue
		}

		trackedFlights = append(trackedFlights, data.ident)
		if data.ident == loc.Flight {
			address = addr
			break
		}
	}
	if address == 0 {
		return fmt.Errorf("flight %s not currently tracked (%s)", loc.Flight, strings.Join(trackedFlights, ", "))
	}

	lp.name = loc.Flight
	lp.focusAddress = address
	return nil
}

func (lp *sbsLocationProvider) Location() (math.LLACoords, time.Time, string, bool) {
	pd, ok := lp.dataCache[lp.focusAddress]
	if !ok {
		return math.LLACoords{}, "Flight " + lp.name, false
	}

	emergency := ""
	if pd.emergency != "" {
		emergency = ", indicating " + pd.emergency
	}

	return pd.loc, pd.updatedAt, "Flight " + lp.name + emergency, true
}

type planeData struct {
	ident              string
	loc                math.LLACoords
	updatedAt          time.Time
	lastWrittenWasEven bool
	evenWritten        bool
	evenCPRLat         uint32
	evenCPRLon         uint32
	oddWritten         bool
	oddCPRLat          uint32
	oddCPRLon          uint32

	emergency string
}

func (lp *sbsLocationProvider) connectDecoder(host string) error {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(conn)
	go func() {
		for {
			for scanner.Scan() {
				// Remove leading *
				rawLine := scanner.Text()
				lineLen := len(rawLine)
				if rawLine[0:1] != "*" || rawLine[lineLen-1:lineLen] != ";" {
					continue
				}

				binaryData, err := hex.DecodeString(rawLine[1 : lineLen-1])
				if err != nil {
					log.Println(err)
					continue
				}

				mms, err := modeSReader.ReadMessage(binaryData)
				if err != nil {
					log.Println(err)
					continue
				}
				if mms.GetDownLinkFormat() != 17 {
					continue
				}
				// It must be a type 17 at this point
				mess17 := mms.(*modeSMessages.MessageDF17)
				messageADSB, _, err := adsbReader.ReadADSBMessage(adsb.ReaderLevel2, false, false, mess17.MessageExtendedSquitter)

				address := mess17.AddressAnnounced.Address

				pd := planeData{}
				if data, ok := lp.dataCache[address]; ok {
					pd = data
				}

				switch typedMessage := messageADSB.(type) {
				case *messages.Format11V2:
					pd = updateLocation(pd, typedMessage.Altitude.AltitudeInFeet, typedMessage.CPRFormat, typedMessage.EncodedLatitude, typedMessage.EncodedLongitude)
				case *messages2.Format04:
					pd = updateCallsign(pd, typedMessage)
				case *messages3.Format19GroundSpeedNormal:
				case *messages4.Format31AirborneV2:
				case *messages5.Format29Subtype1:
				case *messages6.Format28StatusV2:
					pd = updateEmergency(pd, typedMessage.EmergencyPriorityStatus)
				case *messages.Format13V2:
					pd = updateLocation(pd, typedMessage.Altitude.AltitudeInFeet, typedMessage.CPRFormat, typedMessage.EncodedLatitude, typedMessage.EncodedLongitude)
				case *messages.Format12V2:
					pd = updateLocation(pd, typedMessage.Altitude.AltitudeInFeet, typedMessage.CPRFormat, typedMessage.EncodedLatitude, typedMessage.EncodedLongitude)
				default:
					fmt.Printf("-- %T --\n%s\n", typedMessage, typedMessage.ToString())
				}

				lp.dataCache[address] = pd
			}
			if err := scanner.Err(); err != nil {
				log.Println("🛑 Provider: ADS-B no longer available, server communication interrupted", err)
			}
		}

	}()
	return nil
}

func updateLocation(pd planeData, alt int, cpr fields.CompactPositionReportingFormat, cprLat fields.EncodedLatitude, cprLon fields.EncodedLongitude) planeData {
	pd.loc.Altitude = math.Meters(float64(alt) * 0.3048)

	if cpr == fields.CPRFormatEven {
		pd.evenCPRLat = uint32(cprLat)
		pd.evenCPRLon = uint32(cprLon)
		pd.evenWritten = true
		pd.lastWrittenWasEven = true
	} else {
		pd.oddCPRLat = uint32(cprLat)
		pd.oddCPRLon = uint32(cprLon)
		pd.oddWritten = true
		pd.lastWrittenWasEven = false
	}

	if pd.evenWritten && pd.oddWritten {
		lat, lon, err := geo.GetCPRExactPosition(pd.evenCPRLat, pd.evenCPRLon, pd.oddCPRLat, pd.oddCPRLon, pd.lastWrittenWasEven)
		if err == nil {
			pd.loc.Latitude = math.Degrees(lat)
			pd.loc.Longitude = math.Degrees(lon)
			pd.updatedAt = time.Now()
		}
	}
	return pd
}

func updateCallsign(pd planeData, m4 *messages2.Format04) planeData {
	ident := strings.Trim(string(m4.AircraftIdentification), " ")
	pd.ident = ident
	return pd
}

func updateEmergency(pd planeData, status fields2.EmergencyPriorityStatus) planeData {
	switch status {
	case fields2.EPSNoEmergency:
		pd.emergency = ""
	case fields2.EPSGeneralEmergency:
		pd.emergency = "a general emergency"
	case fields2.EPSLifeguardMedical:
		pd.emergency = "a medical emergency"
	case fields2.EPSMinimumFuel:
		pd.emergency = "a fuel emergency"
	case fields2.EPSNoCommunication:
		pd.emergency = "a communications emergency"
	case fields2.EPSUnlawfulInterference:
		pd.emergency = "an unlawful interference emergency"
	case fields2.EPSDownedAircraft:
		pd.emergency = "a crash"
	}
	return pd
}
