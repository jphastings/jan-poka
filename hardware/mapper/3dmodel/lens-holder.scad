use <../../../vendor/openscad/hardware/screws.scad>;

module lens(diameter, f, max_thick) {
  x_shift = f-max_thick/2;
  intersection() {
    translate([x_shift, 0, 0])
      sphere(r=f);
    translate([-x_shift, 0, 0])
      sphere(r=f);
  translate([-max_thick/2,0,0])
      rotate([0, 90, 0])
      cylinder(max_thick, d=diameter);
  }
}

module holder(r, h, thick = 10, countersunk = false) {
  joint = 3;

  module arm() {
    hole_t = [0, 0, r];
    hole_br = [0, r * sin(120), r * cos(120)];
    hole_bl = [0, r * sin(240), r * cos(240)];
    module arm_piece() {
      hull() {
        translate(hole_t)
          rotate([0, 90, 0])
          cylinder(h, d=thick);
        translate(hole_br)
          rotate([0, 90, 0])
          cylinder(h, d=thick);
      }
    }

    module thread() {
      hull() {
        translate(hole_t)
          rotate([0, 90, 0])
          cylinder(0.5, d=joint+2);
        translate([0, 0, r+thick/2]) cube(0.5);
      }
    }

    module joint() {
      rotate([0, 90, 0])
      if (countersunk) {
          countersunk_wallmount(joint, h);
        } else {
          countersunk_hex(5.5, 2, joint, h);
        }
    }

    difference() {
      union() {
        arm_piece();
        rotate([120, 0, 0]) arm_piece();
        rotate([240, 0, 0]) arm_piece();
      }

      lens(90.5, 200, 13.6);

      translate(hole_t)
        joint();
      translate(hole_br)
        joint();
      translate(hole_bl)
        joint();

      rotate([120, 0, 0]) thread();
      rotate([240, 0, 0]) thread();
    }
  }

  arm();
}

detail = $preview ? 64 : 256;

rotate([0, -90, 0])
  holder(60, 6, countersunk = false, $fn = detail);
rotate([180, -90, 0]) translate([0, 90, -30])
  holder(60, 6, countersunk = true, $fn = detail);
