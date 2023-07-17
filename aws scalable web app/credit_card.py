import re
ccn = re.compile(r'[4-6][0-9]{3}[-]?[0-9]{4}[-]?[0-9]{4}[-]?[0-9]{4}')

def check_consecutive(s):
    n = len(s)
    ch = s[0]
    count = 1
    for i in xrange(1, n):
        if s[i] == '-':
            continue
        if s[i] == ch:
            count += 1
            if count == 4:
                return False
        else:
            ch = s[i]
            count = 1
    return True

def check_hyphen(s):
    if s.count('-') == 0 or s.count('-') == 3:
        return True
    else:
        return False

N = input()
for _ in xrange(N):
    s = raw_input().strip()
    if ccn.match(s) and check_consecutive(s) and check_hyphen(s):
        print 'Valid'
    else:
        print 'Invalid'