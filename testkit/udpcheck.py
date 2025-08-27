

import socket
import socks  # pip3 install pysocks
import struct
import random

UDP_ECHO_TEST_SERVER_HOST = 's3.mojotv.cn'
UDP_ECHO_TEST_SERVER_PORT = 8080

V2RAY_SOCKS5_HOST = '127.0.0.1'
V2RAY_SOCKS5_PORT = 1088
V2RAY_USERNAME = None  # Set to username if authentication is required
V2RAY_PASSWORD = None  # Set to password if authentication is required


def dns_query_over_socks5():
    # socks5 server support udp dns query
    # domain baidu.com nameserver 114.114.114.114 53
    
    # Create SOCKS5 UDP socket
    s = socks.socksocket(socket.AF_INET, socket.SOCK_DGRAM)
    try:
        s.set_proxy(socks.SOCKS5, V2RAY_SOCKS5_HOST, V2RAY_SOCKS5_PORT, False, V2RAY_USERNAME, V2RAY_PASSWORD)
        
        # Construct DNS query for baidu.com (A record)
        transaction_id = random.randint(0, 65535)
        flags = 0x0100  # Standard query, recursion desired
        questions = 1
        answers = 0
        authority = 0
        additional = 0
        
        # DNS header (12 bytes)
        header = struct.pack('!HHHHHH', 
                           transaction_id, flags, questions, 
                           answers, authority, additional)
        
        # Question section for baidu.com
        domain = 'baidu.com'
        qname = b''
        for part in domain.split('.'):
            qname += bytes([len(part)]) + part.encode()
        qname += b'\x00'  # Null terminator
        qtype = 1  # A record
        qclass = 1  # IN (Internet)
        question = qname + struct.pack('!HH', qtype, qclass)
        
        dns_query = header + question
        
        # Send DNS query to 114.114.114.114:53
        s.sendto(dns_query, ('114.114.114.114', 53))
        
        # Receive response
        response, address = s.recvfrom(4096)
        
        if len(response) < 12:
            print("Invalid DNS response")
            return
        
        # Parse DNS response header
        resp_id, resp_flags, resp_questions, resp_answers, resp_authority, resp_additional = struct.unpack('!HHHHHH', response[:12])
        
        print(f"DNS Query ID: {transaction_id}")
        print(f"Response ID: {resp_id}")
        print(f"Response flags: {hex(resp_flags)}")
        print(f"Questions: {resp_questions}, Answers: {resp_answers}")

        # print the domain's dns record IP (human readable format)
        if resp_answers > 0:
            # Skip question section to get to answer section
            offset = 12
            # Skip question section (domain + 4 bytes)
            while response[offset] != 0:
                offset += 1
            offset += 5  # null byte + qtype(2) + qclass(2)
            # Parse first answer (skip name(2), type(2), class(2), ttl(4), rdlength(2))
            offset += 2 + 2 + 2 + 4 + 2
            ip_bytes = response[offset:offset+4]
            ip_str = ".".join(str(b) for b in ip_bytes)
            print(f"DNS Record IP: {ip_str}")
        else:
            print("No DNS answer section found")

        if resp_answers > 0:
            print("DNS query successful - received answers")
        else:
            print("DNS query failed - no answers received")
            
    except socks.ProxyError as e:
        print(f"Proxy error: {e.msg}")
    except socket.error as e:
        print(f"Socket error: {repr(e)}")
    except Exception as e:
        print(f"Error: {repr(e)}")
    finally:
        s.close()


def test_udp_echo():
    s = socks.socksocket(
        socket.AF_INET, socket.SOCK_DGRAM
    )  # Same API as socket.socket in the standard lib
    try:
        s.set_proxy(
            socks.SOCKS5, V2RAY_SOCKS5_HOST, V2RAY_SOCKS5_PORT, False, V2RAY_USERNAME, V2RAY_PASSWORD
        )  # SOCKS4 and SOCKS5 use port 1080 by default
        # Can be treated identical to a regular socket object
        # Raw DNS request
        req = '1234'.encode()
        s.sendto(req, (UDP_ECHO_TEST_SERVER_HOST, UDP_ECHO_TEST_SERVER_PORT))
        (rsp, address) = s.recvfrom(4096)
        # check req and rsp equality
        if len(rsp) < 2:
            print("Invalid response")
        if rsp[0] == req[0] and rsp[1] == req[1] and len(rsp) == len(req):
            print("UDP echo server test")
            # print req as hex string
            print("req: " + " ".join("{:02x}".format(c) for c in req))
            print("res: " + " ".join("{:02x}".format(c) for c in rsp))

        else:
            print("Invalid response")
    except socks.ProxyError as e:
        print(e.msg)
    except socket.error as e:
        print(repr(e))



    

if __name__ == "__main__":
    # dns_query_over_socks5()
    # print("DNS query completed")
    print("-----------------")
    test_udp_echo()