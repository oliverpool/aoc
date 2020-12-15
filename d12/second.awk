#!/usr/bin/awk -f


function abs(v) {return v < 0 ? -v : v}

BEGIN {
    wx = 10
    wy = 1
    sx = 0
    sy = 0
    PI = 3.1415926535897932384626433
}
/^N/ {
    wy += substr($0, 2)
}
/^S/ {
    wy -= substr($0, 2)
}

/^E/ {
    wx += substr($0, 2)
}
/^W/ {
    wx -= substr($0, 2)
}

/^F/ {
    sx += substr($0, 2) * wx
    sy += substr($0, 2) * wy
}
/^L/ {
    angle = substr($0, 2)*PI/180.0
    wwy = wy
    wwx = wx
    wx = wwx*cos(angle) - wwy*sin(angle)
    wy = wwx*sin(angle) + wwy*cos(angle)
}
/^R/ {
    angle = -substr($0, 2)*PI/180.0
    wwy = wy
    wwx = wx
    wx = wwx*cos(angle) - wwy*sin(angle)
    wy = wwx*sin(angle) + wwy*cos(angle)
}

/-/ {
    print $0
    print "Sum"
    print sx "," sy
    print wx "," wy
    print ""
}

END {
    print abs(sx) + abs(sy)
}
