def calc_hamming_code(d: str):
    data=list(d).copy()
    data.reverse()
    c,ch,j,r,h=0,0,0,0,[]

    while ((len(d)+r+1)>(pow(2,r))):
        r=r+1

    for i in range(0,(r+len(data))):
        p=(2**c)

        if(p==(i+1)):
            h.append(0)
            c=c+1

        else:
            h.append(int(data[j]))
            j=j+1

    for parity in range(0,(len(h))):
        ph=(2**ch)
        if(ph==(parity+1)):
            startIndex=ph-1
            i=startIndex
            toXor=[]

            while(i<len(h)):
                block=h[i:i+ph]
                toXor.extend(block)
                i+=2*ph

            for z in range(1,len(toXor)):
                h[startIndex]=h[startIndex]^toXor[z]
            ch+=1

    h.reverse()
    # print('Hamming code generated would be:- ', end="")
    return ''.join(map(str, h))


def detect_error(d: str):
    data=list(d).copy()
    data.reverse()
    c,ch,j,r,error,h,parity_list,h_copy=0,0,0,0,0,[],[],[]

    for k in range(0,len(data)):
        p=(2**c)
        h.append(int(data[k]))
        h_copy.append(data[k])
        if(p==(k+1)):
            c=c+1
            
    for parity in range(0,(len(h))):
        ph=(2**ch)
        if(ph==(parity+1)):

            startIndex=ph-1
            i=startIndex
            toXor=[]

            while(i<len(h)):
                block=h[i:i+ph]
                toXor.extend(block)
                i+=2*ph

            for z in range(1,len(toXor)):
                h[startIndex]=h[startIndex]^toXor[z]
            parity_list.append(h[parity])
            ch+=1
    parity_list.reverse()
    error=sum(int(parity_list) * (2 ** i) for i, parity_list in enumerate(parity_list[::-1]))
    return error

def correct_hamming_code(d: str, error: int):
    data=list(d).copy()
    data.reverse()

    if((error)==0):
        return d

    elif((error)>=len(data)):
        return False

    else:
        print('Error is in',error,'bit')

        data[error-1] = ('1' if data[error-1]=='0' else '0')
        
        # print('After correction hamming code is:- ')
        data.reverse()
        return ''.join(map(str, data))

if __name__ == "__main__":
    h1 = calc_hamming_code('1011001')
    print(h1) # 10101001110

    h2 = '11101001110'
    h2 = correct_hamming_code(h2, detect_error(h2))
    assert h1 == h2