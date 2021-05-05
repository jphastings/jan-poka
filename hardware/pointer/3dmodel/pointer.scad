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

outer_stepper_offset = nema_width + 5;
bearing_support_w = 5.8;
bearing_support_h = 2;
bearing_d = 10.75;
bearing_h = 5;
axel_spacing = (bearing_d - axel_d) / 4;

belt_w = 6;
belt_idler_w = 1;
assumed_max_gear_diameter = 80;
outer_gear_height = nema_shaft_length - min_grubbable - bearing_support_h - 1;
outer_gear_timing_belt_d = 35;
timing_teeth = 26;

pointer_extension = 5;

joint_breadth = bearing_h*3;
joint_dia = pointer_extension+bearing_d;
joint_edgevec = [-joint_breadth/2, 0, nema_shaft_length + pointer_extension + joint_dia/2];


linear_extrude(1) stepper_motor_mount(17, 0, true, tolerance);
translate([0, outer_stepper_offset, bearing_support_h - nema_shaft_length + belt_w])
  linear_extrude(1)
  stepper_motor_mount(17, motor_slide, true, tolerance);

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

joint();

// Eyeballed — due to the gear projecting below the horizon
// h_offset = 0.7813;
h_offset = 0;

teeth = 48;
// Eyeballed
gear_d = 59.2;
// Eyeballed, to get the given height
gear_m = gear_d/teeth;
// This is guessed for the height needed
//gear_w = 17.39;
gear_w = 10;


module horizontal_gear() {
  translate([0, 0, bearing_support_h]) difference() {
    union() {
      difference() {
        translate([0, 0, gear_m*h_offset])
          bevel_pair(n1=teeth, n2=teeth, w=gear_w, m=gear_m, only=1);

        // Space for pulley
        //bearing(outer_gear_timing_belt_d, assumed_max_gear_diameter, belt_w+belt_idler_w);
      }

      //GT2(timing_teeth, belt_w, belt_idler_w);
    }

    // Shaft hole
    cylinder(outer_gear_height, d=(axel_d+bearing_d)/2);
    // Lower bearing
    cylinder(bearing_h+tolerance, d=bearing_d+2*tolerance);
    // Upper bearing
    translate([0, 0, outer_gear_height - bearing_h])
      cylinder(bearing_h+tolerance, d=bearing_d+2*tolerance);

  }
}

horizontal_gear();

module vertical_driving_gear() {
  axel_hole_d = axel_d+2*tolerance;

  difference() {
    union() {
      translate([0, 0, gear_m*h_offset])
        bevel_pair(n1=teeth, n2=teeth, w=gear_w, m=gear_m, only=2);

      translate([0, 0, outer_gear_height])
        cylinder(min_grubbable, d = axel_hole_d+2*grub_thickness);
    }

    cylinder(outer_gear_height + min_grubbable, d=axel_hole_d);
    translate([0, 0, outer_gear_height + bearing_h/2])
        // 180/16 cos there are 16 teeth and we want to rotate the grub axel_hole_d
        // half a tooth, so the allen key can fit through a gap
        rotate([90, 0, 180/teeth])
        cylinder(3*grub_thickness, d=grub_d);
  }
}

translate([-32.5, 0, 34+h_offset])
  rotate([0, 90, 0])
  rotate([0, 0, 180/teeth])
  vertical_driving_gear();

module vertical_guiding_gear() {
  difference() {
    translate([0, 0, gear_m*h_offset])
      bevel_pair(n1=teeth, n2=teeth, w=gear_w, m=gear_m, only=2);

    bearing(axel_d, bearing_d+2*tolerance, bearing_h);
    translate([0, 0, outer_gear_height - bearing_h])
      bearing(axel_d, bearing_d+2*tolerance, bearing_h);

    cylinder(outer_gear_height, d=(axel_d+bearing_d)/2);
  }
}

// vertical_guiding_gear();
translate([-34.5-50, 0, 34.5])
  rotate([0, 90, 0])
  %rod(5, 200);
