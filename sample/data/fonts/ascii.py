from sys import stdout

for i in range(16):
    for j in range(16):
        n = i*16+j
        if 32 <= n < 128:
            c = chr(n)
        else:
            c = '.'
        stdout.write(c)
    stdout.write('\n')
