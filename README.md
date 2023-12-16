# 1nject0r
Tools for specific types of Header Injections

# Author 
 adrianalvird [ https://adrianalvird.vercel.app ]

# HHI
 Host Header Injection is a type of Injection where host header can be manupulated .
 Usage :
 
    cat urls.txt | hhi -s <your_site> -m <request_type>
    cat urls.txt | hhi -s mysite.com
    cat urls.txt | hhi -s burp-collaborator.net -m POST
