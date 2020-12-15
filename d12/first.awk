#!/usr/bin/awk -f


function abs(v) {return v < 0 ? -v : v}

BEGIN {
    vertical = 0
    horizontal = 0
    direction = 0
    PI = 3.1415926535897932384626433
}
/^N/ {
    vertical += substr($0, 2)
}
/^S/ {
    vertical -= substr($0, 2)
}

/^E/ {
    horizontal += substr($0, 2)
}
/^W/ {
    horizontal -= substr($0, 2)
}

/^F/ {
    horizontal += substr($0, 2) * cos(direction*PI/180.0)
    vertical += substr($0, 2) * sin(direction*PI/180.0)
}
/^L/ {
    direction += substr($0, 2)
}
/^R/ {
    direction -= substr($0, 2)
}

/-/ {
    print $0
    print "Sum"
    print horizontal
    print vertical
    print ""
}

END {
    print abs(vertical) + abs(horizontal)
}
