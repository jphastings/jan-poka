use <../../../vendor/openscad/hardware/rods.scad>;
use <../../../vendor/openscad/hardware/screws.scad>;
use <../../../vendor/openscad/MCAD/motors.scad>;
use <../../../vendor/openscad/MCAD/stepper.scad>;

$fn=25;
tolerance = 0.1;
motor_slide = 5;
nema_depth = 60;
nema_width = 42.3;
nema_shaft_length = 25.5;
axel_d = 5;

min_grubbable = 5;
grub_thickness = 3;

outer_stepper_offset = nema_width + 5;
bearing_support_w = 5.8;
bearing_support_h = min_grubbable;
bearing_d = 10.75;
bearing_h = 5;
axel_spacing = (bearing_d - axel_d) / 4;
belt_w = 6.8;

outer_gear_height = 15;
outer_gear_timing_belt_d = 30;

pointer_extension = 5;


stepper_motor_mount(17, 0, true, tolerance);
translate([0, outer_stepper_offset, bearing_support_h - nema_shaft_length + belt_w])
  stepper_motor_mount(17, motor_slide, true, tolerance);

module pointer() {
  rotate_extrude()
    polygon([
      [0, nema_shaft_length + pointer_extension],
      [0, nema_shaft_length + tolerance],
      [axel_d/2 + tolerance, nema_shaft_length + tolerance],
      [axel_d/2 + tolerance, nema_shaft_length - min_grubbable],
      [axel_d/2 + grub_thickness, nema_shaft_length - min_grubbable],
      [axel_d/2 + grub_thickness, nema_shaft_length + pointer_extension],
    ]);

}

pointer();

module outer_gear() {
  translate([0, 0, bearing_support_h])
    %bearing(axel_d, bearing_d, bearing_h);

  translate([0, 0, bearing_support_h + outer_gear_height - bearing_h])
    %bearing(axel_d, bearing_d, bearing_h);

  rotate([90, 0, 90])
  // rotate_extrude()
    polygon([
      [bearing_d/2+tolerance, bearing_support_h + outer_gear_height],
      [bearing_d/2+tolerance, bearing_support_h + outer_gear_height],
      [bearing_d/2+tolerance, bearing_support_h + outer_gear_height - bearing_h - tolerance],
      [axel_d/2+axel_spacing, bearing_support_h + outer_gear_height - bearing_h - tolerance],
      [axel_d/2+axel_spacing, bearing_h + bearing_support_h + tolerance],
      [bearing_d/2+tolerance, bearing_h + bearing_support_h + tolerance],
      [bearing_d/2+tolerance, bearing_support_h],
      [(bearing_d + outer_gear_timing_belt_d)/2, bearing_support_h],
      [(bearing_d + outer_gear_timing_belt_d)/2, bearing_support_h + belt_w + tolerance]
    ]);
}