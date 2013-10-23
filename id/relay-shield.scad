xdim = 30;
ydim = 40;
zdim = 25;
wing = 20;
slot = 10;
wall = 3;
screw = 4;

module wing() {
	difference() {
		cube([xdim, wing, wall]);
		translate([xdim / 2, wing / 2, 0]) {
			cylinder(r=screw, h=wall);
		}
	}
}

union() {
	cube([wall, ydim, zdim]);
	difference() {
		union() {
			cube([xdim, wall, zdim]);
			translate([0, ydim - wall, 0]) {
				cube([xdim, wall, zdim]);
			}
		}
		translate([2*wall, 0, 2*wall]) {
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