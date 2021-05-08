use <../../../vendor/openscad/hardware/rods.scad>;
use <../../../vendor/openscad/hardware/screws.scad>;
use <../../../vendor/openscad/hardware/gt2.scad>;
use <../../../vendor/openscad/MCAD/motors.scad>;
use <../../../vendor/openscad/MCAD/stepper.scad>;
use <../../../vendor/openscad/PolyGear/PolyGear.scad>;

$fn=25;
if (!$preview) {
  $fn=100;
}


tolerance = 0.1;
motor_slide = 5;
nema_depth = 60;
nema_width = 42.3;
nema_shaft_length = 21.5;
axel_d = 5;

min_grubbable = 5;
grub_thickness = 3;
grub_d = 3;

outer_stepper_offset = nema_width + 2;
bearing_support_w = 5.8;
bearing_support_h = 2;
bearing_d = 10.75;
bearing_h = 5;
axel_spacing = (bearing_d - axel_d) / 4;

timing_pulley_height = (6.8 + 1.1);
timing_pulley_top_d = 15.8;

pointer_extension = 5;

joint_breadth = bearing_h*3;
joint_dia = pointer_extension+bearing_d;
joint_edgevec = [-joint_breadth/2, 0, nema_shaft_length + pointer_extension + joint_dia/2];

module drivers() {
  wall_thickness = 2;

  linear_extrude(wall_thickness) stepper_motor_mount(17, 0, true, tolerance);

  translate([0, outer_stepper_offset, timing_pulley_height - nema_shaft_length])
    linear_extrude(wall_thickness)
    stepper_motor_mount(17, motor_slide, true, tolerance);

  translate([0, outer_stepper_offset, 0])
    %bearing(axel_d, timing_pulley_top_d, timing_pulley_height);
}

module joint() {
  difference() {
    hull() {
      translate([0, 0, nema_shaft_length - min_grubbable])
        cylinder(min_grubbable + pointer_extension, d=2*grub_thickness + axel_d);

      translate(joint_edgevec)
        rotate([0, 90, 0])
        cylinder(joint_breadth, d=joint_dia);
    }


    translate(joint_edgevec) {
      rotate([0, 90, 0])
        bearing(axel_d, bearing_d+2*tolerance, bearing_h);

      translate([joint_breadth-bearing_h, 0, 0])
        rotate([0, 90, 0])
        bearing(axel_d, bearing_d+2*tolerance, bearing_h);

      rotate([0, 90, 0])
        cylinder(joint_breadth, d=(axel_d+bearing_d)/2);
    }

    rod(axel_d+2*tolerance, nema_shaft_length);

    translate([0, 0, nema_shaft_length - min_grubbable/2])
      rotate([90, 0, 0])
      cylinder(3*grub_thickness, d=grub_d);
  }
}

module bevel(horizontal = true, position = true) {
  // fudges, cos I can't get the maths right
  ugh = 0.85;
  z_offset = 0.4;

  // constants
  cone_angle = 45;
  n = 32;
  gear_height = 6.8;
  // Calculations
  base_to_center = joint_edgevec[2] - timing_pulley_height;
  H0 = base_to_center / (1 + 2*tan(cone_angle)*ugh / n);
  r0 = H0*tan(cone_angle);
  m = 2*r0 / n;

  tr = (position) ?
    (horizontal) ? [0,0,timing_pulley_height] : [-H0,0,joint_edgevec[2]]
    : [0,0,0];
  ro = (horizontal) ? [0, 0, 0] : [0, 90, 0];

  translate(tr)
    rotate(ro)
    difference() {
      union() {
        translate([0, 0, m * ugh])
          bevel_gear(cone_angle=cone_angle, z=(gear_height-z_offset), m=m, n=n, lift=false);

        if (!horizontal) {
          translate([0, 0, gear_height])
            cylinder(min_grubbable, d=axel_d + grub_thickness*2);

          %rod(5, 200);
        }
      }

      if (horizontal) {
        cylinder(gear_height, d = timing_pulley_top_d + 2*tolerance);
      } else {
        cylinder(gear_height+ min_grubbable, d = axel_d + 2*tolerance);

        translate([0, 0, gear_height + min_grubbable/2])
          rotate([90, 0, 0])
          cylinder(3*grub_thickness, d=grub_d);
      }
    }

}

drivers();
joint();
bevel(horizontal = true);
bevel(horizontal = false);
