import io
import os
import re

def buildProto(pbpath,outpath,pbfile):
    cmdstr = 'protoc.exe --go_out={0} -I{1} {1}/{2}.proto'.format(outpath,pbpath,pbfile)
    os.system(cmdstr)
    file_data = ""
    tags = {}
    tagsstr = ""
    file = "{0}/{1}.pb.go".format(outpath,pbfile)
    with io.open(file, "r", encoding='utf-8') as f:
        for line in f:
            if 'tags:' in line:
                for v in re.findall(r"`(.+?)`",line)[0].split(' '):
                    tag = v.split(':')
                    tags[tag[0]] = tag[1]
                for v in re.findall(r"tags:{(.+?)}",line)[0].split(' '):
                    tag = v.split(':')
                    tags[tag[0]] = tag[1]
                for key,value in tags.items():
                    tagsstr += "{0}:{1} ".format(key,value)
                line = re.sub(r"`(.+?)`", "`{0}`".format(tagsstr[0:len(tagsstr)-1]), line)
            file_data += line
    with io.open(file,"w",encoding='utf-8') as f:
        f.write(file_data)


buildProto('./core','./core','proto')

