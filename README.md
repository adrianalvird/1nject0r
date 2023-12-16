# 1nject0r
Tools for specific types of Header Injectors
     Author  :  adrianalvird [ https://adrianalvird.vercel.app ]

[ For the customization :

    1. copy the script like hhi.go 
    2. then modify as you needed
    3. go build hhi.go 
    4. sudo cp hhi /usr/bin

]


# HHI
 Host Header Injection is a type of Injection where host header can be manupulated .
 Usage :
 
    cat urls.txt | hhi -s <your_site> -m <request_type>
    cat urls.txt | hhi -s mysite.com
    cat urls.txt | hhi -s burp-collaborator.net -m POST
    cat urls.txt | hhi -s attackersite.com | grep [+]     //advance usage

# CRLF
CRLF Injector Tool added inside Injector 
Usage: 

    cat urls.txt | crlf -mc 200                //mc = specific status code 
    cat urls.txt | grep [+]     //advance usage  
    
    

#
