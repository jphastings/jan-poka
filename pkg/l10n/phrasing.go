package l10n

import (
	"fmt"
	. "github.com/jphastings/corviator/pkg/math"
	"math"
)

const exactAngleTolerance = Degrees(5)

var compassPoints = []string{
	"North", "North North East", "North East", "East North East",
	"East", "East South East", "South East", "South South East",
	"South", "South South West", "South West", "West South West",
	"West", "West North West", "North West", "North North West"}

var aboveBelow = map[bool]string{true: "above", false: "below"}
var upDown = map[bool]string{true: "up", false: "down"}
var conversationalElevation = map[bool]string{true: ", and %d degrees %s the horizon", false: " %d degrees %s"}
var conversationalDir = map[bool]map[bool]string{true: aboveBelow, false: upDown}
var conversationalAhead = map[bool]string{true: "", false: " dead ahead"}

func Phrase(name string, bearing AERCoords, isFirstTrack bool) string {
	if isFirstTrack {
		return fmt.Sprintf("Turn to face%s, and look%s. %.0fkm that way, you'll find %s.",
			compassHeading(bearing.Azimuth),
			elevation(bearing.Elevation, false),
			bearing.Range/1000,
			name)
	} else {
		return fmt.Sprintf("%s is now%s%s, %.0fkm away.",
			name,
			compassHeading(bearing.Azimuth),
			elevation(bearing.Elevation, true),
			bearing.Range/1000)
	}
}

func compassHeading(azimuth Degrees) string {
	approxDir := int((azimuth + 11.25) / 22.5)
	compassPoint := compassPoints[approxDir]

	accuracy := ModDeg(Degrees(math.Abs(float64(approxDir)*22.5 - float64(azimuth))))
	if accuracy < exactAngleTolerance {
		return " " + compassPoint
	}

	return " roughly " + compassPoint
}

func elevation(elevation Degrees, conversational bool) string {
	if elevation > 90-exactAngleTolerance {
		return " straight up"
	}
	if elevation < -90+exactAngleTolerance {
		return " straight down"
	}
	if elevation > -exactAngleTolerance && elevation < exactAngleTolerance {
		return conversationalAhead[conversational]
	}

	roughAngle := int(exactAngleTolerance) * int(math.Abs(float64(elevation))/float64(exactAngleTolerance))
	return fmt.Sprintf(conversationalElevation[conversational],
		roughAngle,
		conversationalDir[conversational][elevation > 0])
}
