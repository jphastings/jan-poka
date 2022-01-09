module countersunk_wallmount(screw_d, height, countersink_angle=90, tolerance=0.1) {
  sink_h = screw_d;
  sink_d = 2 * sink_h * tan(countersink_angle/2);
  translate([0, 0, height - sink_h]) {
    cylinder(sink_h, d1=2*tolerance, d2=sink_d+2*tolerance);
  }
    cylinder(height, d=screw_d+2*tolerance);
}

module countersunk_hex(hex_small_d, hex_height, screw_d, height, tolerance = 0.1) {
  cylinder(height, d=screw_d+2*tolerance);

  hex_d = hex_small_d/cos(30);
  chamfer_h = hex_d/4;
  translate([0, 0, height - hex_height]) cylinder(hex_height, d=hex_d, $fn=6);
  translate([0, 0, height - hex_height - chamfer_h]) cylinder(chamfer_h, d1=0, d2=hex_d, $fn=6);
}