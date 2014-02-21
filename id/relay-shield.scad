wall = 2;
xdim = 30;
ydim = 50 + 2*wall;
zdim = 28;
wing = 20;
slot = 8;
screw = 2.5;
$fn = 6;

module wing() {
	difference() {
		cube([xdim, wing, wall]);
		translate([xdim / 2, wing / 2, 0]) {
			cylinder(r=screw, h=wall);
		}
	}
}

module shield() {
	union() {
		cube([wall, ydim, zdim]);
		difference() {
			union() {
				cube([xdim, wall, zdim]);
				translate([0, ydim - wall, 0]) {
					cube([xdim, wall, zdim]);
				}
			}
			translate([2*wall, 0, zdim - 3*wall - slot]) {
				cube([xdim, ydim + 2*wall, slot]);
			}
		}
		translate([0, 0, zdim - wall]) {
			cube([xdim, ydim, wall]);
		}
		translate([0, -wing, 0]) {
			wing();
		}
		translate([0, ydim, 0]) {
			wing();
		}
	}
}

translate([zdim / 2, -wall - ydim / 2, 0]) rotate(270, [0, 1, 0])
  shield();
