#include <algorithm>
#include <fstream>
#include <string>
#include <cstdint>
#include <algorithm>  
#include <mutex>          // std::call_once, std::once_flag
#include <math.h> // round()
#include <stdint.h>
#include "LZJD.h"
#include "MurmurHash3.h"
using namespace std;

#ifdef __cplusplus
extern "C" {
#endif

LZJD::LZJD() {
}

LZJD::~LZJD() {
}


std::vector<int32_t> getAllHashes(std::vector<char>& bytes)
{
    std::vector<int32_t> ints;

    std::unordered_set<int32_t> x_set;
    MurmurHash3 running_hash = MurmurHash3();

    for(char b : bytes) 
    {
        int8_t some_byte = b; //wherever you want to get this
        int32_t hash = running_hash.pushByte(some_byte);

        if (x_set.insert(hash).second)
        {
            //was successfully added, so never seen it before. Put it in! 
            ints.push_back(hash);
            running_hash.reset();
        }

    }

    return ints;
}

std::vector<int32_t> digest(uint64_t k, std::vector<char>& bytes)
{
    std::vector<int32_t> ints = getAllHashes(bytes);

    if(ints.size() > k)
    {
        std::nth_element (ints.begin(), ints.begin()+k, ints.end());
        ints.resize(k);
        std::sort(ints.begin(), ints.end());
    }
    else
    {
        std::sort(ints.begin(), ints.end());
        ints.resize(k);
    }
    
    return ints;
}

int32_t similarity(const std::vector<int32_t>& x_minset, const std::vector<int32_t>& y_minset)
{
    int32_t same = 0;
    vector<int32_t> v3;
    set_intersection(x_minset.begin(),x_minset.end(),
                     y_minset.begin(),y_minset.end(),
                     back_inserter(v3));
    same = v3.size();
    double sim = same / (double) (x_minset.size() + y_minset.size() - same);
    return (int) (round(100*sim));
}

#ifdef __cplusplus
}
#endif