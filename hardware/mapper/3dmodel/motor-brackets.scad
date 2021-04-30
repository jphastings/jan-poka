use <../../../vendor/openscad/hardware/rods.scad>;
use <../../../vendor/openscad/hardware/screws.scad>;
use <../../../vendor/openscad/MCAD/motors.scad>;
use <../../../vendor/openscad/MCAD/stepper.scad>;

$fn=25;

module worm_wall_mount(tolerance = 0.1) {
  wall_thickness = 2;
  nema_17_thickness = 42.20;
  nema_backstep = 18;
  stepper_gap = 4;
  axel_offset = 10;

  bearing_height = 5;
  bearing_od = 13;
  bearing_housing_od = bearing_od+2*wall_thickness;
  bracket_thickness = bearing_height + wall_thickness;

  difference() {
    translate([axel_offset - nema_17_thickness/2, nema_backstep - wall_thickness, 0])
      cube([nema_17_thickness, wall_thickness, nema_17_thickness + 3*wall_thickness]);

    translate([axel_offset, nema_backstep, nema_17_thickness/2 + wall_thickness + stepper_gap])
      rotate([90, 90, 0])
      linear_extrude(wall_thickness)
      stepper_motor_mount(17, 2, true, tolerance);
  }


  translate([0, 0, wall_thickness]) {
    %bearing(4, bearing_od+2*tolerance, bearing_height+tolerance);
    %rod(4, 70);
  }
  translate([0, 0, nema_17_thickness + 4*wall_thickness])
    %bearing(4, bearing_od+2*tolerance, bearing_height+tolerance);

  module bracket() {
    difference() {
      hull() {
        cylinder(bracket_thickness, d=bearing_od+2*wall_thickness);
        translate([axel_offset - nema_17_thickness/2, nema_backstep - wall_thickness, 0])
          cube([nema_17_thickness, wall_thickness, bracket_thickness]);
      }

      translate([0, 0, wall_thickness])
        cylinder(bearing_height+tolerance, d=bearing_od+2*tolerance);
    }
  }

  difference() {
    bracket();
    translate([-5, 12, 0])
      countersunk_wallmount(3.5, bracket_thickness, tolerance = tolerance);
    translate([15, 12, 0])
      countersunk_wallmount(3.5, bracket_thickness, tolerance = tolerance);
    translate([0, 5, 0])
      cylinder(wall_thickness+tolerance, d=3);
    translate([0, -5, 0])
      cylinder(wall_thickness+tolerance, d=3);
  }


  difference() {
    translate([0, 0, nema_17_thickness + 3*wall_thickness])
      bracket();

    rod(9, 70);
    translate([-5, 12, nema_17_thickness + 2*wall_thickness])
      cylinder(bracket_thickness+2*wall_thickness, d=8);
    translate([15, 12, nema_17_thickness + 2*wall_thickness])
      cylinder(bracket_thickness+2*wall_thickness, d=8);
  }
}

translate([20, 0, 0])
  rotate([-90, 0, 0])
  worm_wall_mount();

translate([-20, 0, 0])
  rotate([-90, 0, 0])
  mirror([1, 0, 0]) worm_wall_mount();