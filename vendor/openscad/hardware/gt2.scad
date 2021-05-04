// From https://github.com/rbuckland/openscad.parametric-pulley
module GT2(teeth, pulley_t_ht, idler_ht) {
  include <../parametric-pulley/parametric-pulley-generator.scad>;

  retainer_ht = 0;  // Belt retainer above teeth : height of retainer flange over pulley, standard = 1.5
  retainer_flat_width = 0; // Width of retainer after reaching height above teeth
  idler_flat_width = 0; // Width of idler after reaching height above teeth

  pulley_b_ht = 0;    // pulley base height, standard = 8. Set to 0 idler but no pulley.
  pulley_b_dia = 25;  // pulley base diameter, standard = 20
  no_of_nuts = 0;   // number of captive nuts required, standard = 1
  nut_angle = 90;   // angle between nuts, standard = 90
  nut_shaft_distance = 0; // distance between inner face of nut and shaft, can be negative.


  //  ********************************
  //  ** Scaling tooth for good fit **
  //  ********************************
  /*  To improve fit of belt to pulley, set the following constant. Decrease or increase by 0.1mm at a time. We are modelling the *BELT* tooth here, not the tooth on the pulley. Increasing the number will *decrease* the pulley tooth size. Increasing the tooth width will also scale proportionately the tooth depth, to maintain the shape of the tooth, and increase how far into the pulley the tooth is indented. Can be negative */

  additional_tooth_width = 0.2; //mm

  //  If you need more tooth depth than this provides, adjust the following constant. However, this will cause the shape of the tooth to change.

  additional_tooth_depth = 0.0; //mm

  pulley(
      belt_description          = "GT2 5mm"
      , pulley_OD               = calc_pulley_dia_tooth_spacing (teeth, 5,0.5715) // Set the pulley diameter for a given number of teeth
      , teeth                   = teeth
      , tooth_depth             = 1.969
      , tooth_width             = 3.952
      , additional_tooth_depth  = additional_tooth_depth
      , additional_tooth_width  = additional_tooth_width
      , motor_shaft             = 0 ) tooth_profile_GT2_5mm(height = pulley_t_ht);
}

module tooth_profile_GT2_5mm(height) {
  linear_extrude(height=height+2) polygon([[-1.975908,-0.75],[-1.975908,0],[-1.797959,0.03212],[-1.646634,0.121224],[-1.534534,0.256431],[-1.474258,0.426861],[-1.446911,0.570808],[-1.411774,0.712722],[-1.368964,0.852287],[-1.318597,0.989189],[-1.260788,1.123115],[-1.195654,1.25375],[-1.12331,1.380781],[-1.043869,1.503892],[-0.935264,1.612278],[-0.817959,1.706414],[-0.693181,1.786237],[-0.562151,1.851687],[-0.426095,1.9027],[-0.286235,1.939214],[-0.143795,1.961168],[0,1.9685],[0.143796,1.961168],[0.286235,1.939214],[0.426095,1.9027],[0.562151,1.851687],[0.693181,1.786237],[0.817959,1.706414],[0.935263,1.612278],[1.043869,1.503892],[1.123207,1.380781],[1.195509,1.25375],[1.26065,1.123115],[1.318507,0.989189],[1.368956,0.852287],[1.411872,0.712722],[1.447132,0.570808],[1.474611,0.426861],[1.534583,0.256431],[1.646678,0.121223],[1.798064,0.03212],[1.975908,0],[1.975908,-0.75]]);
}
