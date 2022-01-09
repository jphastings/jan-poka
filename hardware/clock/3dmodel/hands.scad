use <../../../ext/openscad/2d_points/2d_points.scad>;

t = 4;
minuteInnerD = 1.8; // want 1.3
minuteOuterD = 7.5;
minuteLength = 85;
minuteMagnet = 10;
hourInnerD = 4.15; // want 3.7
hourOuterD = 15;
hourLength = minuteLength *2 /3;
hourMagnet = 30;

magnetD = 6; // want 5.8
magnetH = 3; // want 2.6
magnetMass = 0.6; // g

counterweightD = 9.1; // want 8.55
counterweightH = 2.8; // want 2.45
counterweightMass = 1.1; // g
//counterweightMinPos = minuteMagnet * magnetMass / counterweightMass;
counterweightMinPos = 9.8625;
manualOffset = 11; // manually measured an extra 1 magnet was needed at 11mm out
counterweightHourPos = (hourMagnet + manualOffset) * magnetMass / counterweightMass;

function cat(L1, L2) = [for(L=[L1, L2], a=L) a];

module hand(inner, outer, length, widthMultiplier = 1, minWidth = 1.2, tail = false) {
    d = (outer - inner)/2;
    width = widthMultiplier * d;
    pt = inner/2 + 2*d;
    
    shortEnd = [inner/2, 0];
    leftOut = [pt, width/2];
    rightOut = [pt, -width/2];
  
    circleEdge = outer/2 * cos(45);
    
    mid = [rightOut[0]/2, rightOut[1]/2];
      
    //color([0, 0, 255]) translate(mid) cylinder(1.1*t, d=1);
    
    curve = cat([shortEnd], cat(bezier_points([rightOut, mid, [circleEdge, -circleEdge]]), [shortEnd]));
    
    union() {
      difference() {
          cylinder(t, d=outer);
          cylinder(t, d=inner);
      }
    
      linear_extrude(t)
      polygon([
        shortEnd,
        leftOut,
        [length - minWidth/2, minWidth/2],
        [length, 0],
        [length - minWidth/2, -minWidth/2],
        rightOut
      ]);
      
      linear_extrude(t) polygon(points = curve);
      mirror([0, 1, 0]) linear_extrude(t) polygon(points = curve);
      
      if (tail) {
        mirror([1, 0, 0])
        linear_extrude(t)
        polygon([
          [circleEdge, -circleEdge],
          [circleEdge, 0],
          [circleEdge, circleEdge],
          1.5*leftOut,
          1.5*rightOut
        ]);
        echo (1.5*leftOut[0]);
        translate([-pt * 1.5, 0, 0]) {
          cylinder(t, d=1.5*width);
          //translate([0, -0.75*width]) cube(1.5*width);
        }
      }
    }
}

$fn=100;
difference() {
  hand(hourInnerD, hourOuterD, hourLength, 2.5, tail=true);
  translate([hourMagnet, 0, 0])
    cylinder(magnetH, d = magnetD);
  translate([-counterweightHourPos, 0, 0])
    cylinder(counterweightH, d = counterweightD);
}

translate([0, -25, 0]) difference() {
  hand(minuteInnerD, minuteOuterD, minuteLength, 2.8, tail = true);
  translate([minuteMagnet, 0, 0])
    cylinder(magnetH, d = magnetD);
  translate([-counterweightMinPos, 0, 0])
    cylinder(t, d = counterweightD);
}