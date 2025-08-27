

import socket
import socks  # pip3 install pysocks

UDP_ECHO_TEST_SERVER_HOST = 's3.mojotv.cn'
UDP_ECHO_TEST_SERVER_PORT = 8080

V2RAY_SOCKS5_HOST = '127.0.0.1'
V2RAY_SOCKS5_PORT = 1080


if __name__ == "__main__":
    s = socks.socksocket(
        socket.AF_INET, socket.SOCK_DGRAM
    )  # Same API as socket.socket in the standard lib
    try:
        s.set_proxy(
            socks.SOCKS5, V2RAY_SOCKS5_HOST, V2RAY_SOCKS5_PORT, False, user, pwd
        )  # SOCKS4 and SOCKS5 use port 1080 by default
        # Can be treated identical to a regular socket object
        # Raw DNS request
        req = 'abcdexxxx'.encode()
        s.sendto(req, (UDP_ECHO_TEST_SERVER_HOST, UDP_ECHO_TEST_SERVER_PORT))
        (rsp, address) = s.recvfrom(4096)
        # check req and rsp equality
        if len(rsp) < 2:
            print("Invalid response")
        if rsp[0] == req[0] and rsp[1] == req[1] and len(rsp) == len(req):
            print("UDP check passed")
            # print req as hex string
            print("req: " + " ".join("{:02x}".format(c) for c in req))
            print("res: " + " ".join("{:02x}".format(c) for c in rsp))

        else:
            print("Invalid response")
    except socks.ProxyError as e:
        print(e.msg)
    except socket.error as e:
        print(repr(e))
