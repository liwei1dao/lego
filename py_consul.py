#!/usr/bin/env python
 
# encoding: utf-8
import requests
import json
 
host = '10.0.0.9'
port = 8500
 
#consul中微服务名称
serviceNames='dreamfactory'
 
def get_service():
    url = 'http://'+host+':'+str(port)+'/v1/agent/services'
    data = requests.get(url).json()    
    jsonStr = json.dumps(data)
    print('json '+jsonStr)
    keys = []
    for k in data:    
        serviceName=data[k]['Service']
        if serviceName in serviceNames:
             print(data[k]['Service'])
             print(data[k]['ID'])
             keys.append(k)
    return keys
 
def del_service(keys):
    url = 'http://'+host+':'+str(port)+'/v1/agent/service/deregister/'
    for sid in keys:
        requests.put(url+sid)
        print('删除 '+sid)
 
if __name__ == '__main__':
    keys = get_service()
    del_service(keys)