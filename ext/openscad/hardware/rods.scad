// A (ball) bearing with the given inner diameter, outer diameter and height.
//     bearing(4, 8, 4)
module bearing(inner_diameter, outer_diameter, height) {
  rotate_extrude()
    translate([inner_diameter/2, 0])
    square([(outer_diameter - inner_diameter)/2, height]);
}

module rod(outer_diameter, length) {
  cylinder(length, d=outer_diameter);
}