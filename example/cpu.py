#!/usr/bin/env python

import time

print 'start'
start = time.time()
while True:
    if time.time() - start > 30:
        break
