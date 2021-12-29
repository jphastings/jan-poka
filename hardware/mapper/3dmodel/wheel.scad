include <../../../ext/openscad/hardware/bezier.scad>

RAD = 180 / PI;

flange_diameter = 22;
flange_holes = 4;
flange_holes_diameter = 15.5;
flange_hole_diameter = 3.1;
flange_countersink_diameter = 5.5;
flange_countersink_angle = 90;
flange_countersink_h = flange_countersink_diameter / (2 * tan(flange_countersink_angle/2));

cutaway_count = 5;
diameter = 100;
outer_thickness = 8;
rim = 0.5;
inner_thickness = 2;
hole_space = 4;
anchor_d = 2;

assert(anchor_d <= 2*hole_space, "The anchor is too wide to fit on the spokes");

spindle_diameter = (outer_thickness - rim * 2);
end = (diameter - outer_thickness)/2;

thread_diameter = 1.5;


module profile() {
    
    difference() {
        union() {
            translate([end, 0]) square([outer_thickness / 2, outer_thickness/2]);
            square([end, inner_thickness]);
            curve_width = diameter/2 - flange_diameter;
            curve(
                [flange_diameter, inner_thickness],
                [flange_diameter + 0.8 * curve_width, inner_thickness/2],
                [flange_diameter + 0.8 * curve_width, (outer_thickness + inner_thickness)*3/4],
                [diameter/2, outer_thickness]
            );
        }
       
        translate([diameter / 2, outer_thickness / 2]) circle(d = spindle_diameter);
    }
    
    
}


module hole(angle) {
    corner_r = 2;

    txc_y = (corner_r+hole_space);
    txc_x = txc_y * ( 1 / sin(angle) + 1/tan(angle));
    txc_r = sqrt(txc_x*txc_x + txc_y * txc_y);
    
    near = flange_diameter/2 + hole_space - txc_r;
    far = end - hole_space - txc_r;
    half = (near + far) / 2;
    
    assert(near * sin(angle) > 0, "Too many cutaways to fit :(");
    
    union() {
        translate([txc_x, txc_y, 0])
        hull() {
            translate([near * cos(0), near * sin(0), 0]) cylinder(outer_thickness, r = corner_r);
            translate([far * cos(0), far * sin(0), 0]) cylinder(outer_thickness, r = corner_r);
            translate([half * cos(angle/2), half * sin(angle/2), 0]) cylinder(outer_thickness, r = 1);
        }
        
        translate([txc_x, txc_y, 0])
        hull() {
            translate([near * cos(angle), near * sin(angle), 0]) cylinder(outer_thickness, r = corner_r);
            translate([far * cos(angle), far * sin(angle), 0]) cylinder(outer_thickness, r = corner_r);
            translate([half * cos(angle/2), half * sin(angle/2), 0]) cylinder(outer_thickness, r = 1);
        }
        
        
        arc_far_r = sqrt(txc_y^2 + (txc_x+far)^2);
        arc_near_r = sqrt(txc_y^2 + (txc_x+near)^2);
        arc_angle = 1.2; // TODO: Figure this through trig, instead of guessing
        difference() {
            intersection() {
                translate([txc_x, txc_y, 0])
                    linear_extrude(outer_thickness) polygon([
                        [-corner_r/2, -corner_r/2],
                        [(diameter-corner_r/2) * cos(arc_angle), (diameter-corner_r/2) * sin(arc_angle)],
                        [(diameter-corner_r/2) * cos(angle - arc_angle), (diameter-corner_r/2) * sin(angle - arc_angle)]
                    ]);
                cylinder(outer_thickness, r=arc_far_r+corner_r);
            }

            cylinder(outer_thickness, r=arc_near_r - corner_r);
        }
    }
}

module holes(count) {
    angle = 360 / count;
    for(i = [0:count -1]) {
        rotate([0, 0, i * angle])
        hole(angle);
    }
}
    

module curve(start, cp1, cp2, end, steps = $fn) {
    for(step = [steps:1]) {
        polygon(
            points=[
                [end[0], start[1]],
                PointOnBezCubic2D(start, cp1, cp2, end, step/steps),
                PointOnBezCubic2D(start, cp1, cp2, end, (step-1)/steps)],
            paths=[[0,1,2,0]]
        );
    }
}

$fn=100;

if ($preview) { %holes(cutaway_count); } // without difference for preview

difference() {
    rotate_extrude()
        profile();
    if (!$preview) { holes(cutaway_count); }; // this eats CPU, so only on full render

    translate([diameter/2, 0, outer_thickness/2]) rotate([0, -80, 0]) cylinder(spindle_diameter*2, d = thread_diameter);

    for (i = [0:flange_holes - 1]) {
        rotate(360*i/flange_holes) translate([flange_holes_diameter/2, 0, 0]) { cylinder(outer_thickness, d = flange_hole_diameter);
            translate([0, 0, inner_thickness - flange_countersink_h]) cylinder(flange_countersink_h, d1 = 0, d2 = flange_countersink_diameter);
        }
    }
}

//profile();