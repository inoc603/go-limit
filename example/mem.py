#!/usr/bin/env python

import time
s = []
print 'start'
while True:
    s.append(' ' * 128 * 10**6)
    print 'Allocated', len(s) * 128, 'MB'
