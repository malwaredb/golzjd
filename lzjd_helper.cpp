#include <fstream>
#include <iostream>
#include <sstream> 
#include <cstdlib>
#include <cstdint>
#include <string>
#include <cstring>
#include <boost/archive/iterators/base64_from_binary.hpp>
#include <boost/archive/iterators/binary_from_base64.hpp>
#include <boost/archive/iterators/transform_width.hpp>
#include <boost/archive/iterators/ostream_iterator.hpp>
#include <boost/archive/iterators/remove_whitespace.hpp>

#include "LZJD.h"

using namespace std;
namespace bi = boost::archive::iterators;

extern "C" {

const uint64_t digest_size = 1024;

//#define DEBUGHASHES

vector<int32_t> cstring_to_lzjd(char* hash) {
    string line = hash;
    auto first_colon = line.find(":", 0);
    auto second_colon = line.find(":", first_colon + 1);
    string path = line.substr(first_colon + 1, second_colon - first_colon - 1);
    string base64ints = line.substr(second_colon + 1, line.size() - second_colon);
    auto size = base64ints.size();
    while (size > 0 && base64ints[size - 1] == '=')
        size--;
    base64ints = base64ints.substr(0, size);


    //TODO this is not 100% kosher, but C++ is a pain. 

    typedef
    bi::transform_width<
            bi::binary_from_base64<bi::remove_whitespace < string::const_iterator>>,
            8, 6
            >
            base64_dec;

    vector<uint8_t> int_parts;

    copy(
            base64_dec(base64ints.cbegin()),
            base64_dec(base64ints.cend()),
            std::back_inserter(int_parts)
            );

    vector<int32_t> decoded_ints(int_parts.size() / 4);
    for (unsigned int i = 0; i < int_parts.size(); i += 4) {
        //big endian extraction of the right value
        int32_t dec_i = (int_parts[i + 0] << 24) | (int_parts[i + 1] << 16) | (int_parts[i + 2] << 8) | (int_parts[i + 3] << 0);
        decoded_ints[i / 4] = dec_i;
        //                cout << dec_i << ", ";
    }
    return decoded_ints;
}

int32_t lzjd_similarity(char *hash1, char *hash2) {
    try {
        vector<int32_t> l1 = cstring_to_lzjd(hash1);
        vector<int32_t> l2 = cstring_to_lzjd(hash2);
        return similarity(l1, l2);
    } catch(...) {
        return 0;
    }
    return 0;
}

// Functions related to creating LZJD hashes

void readAllBytes(char const* filename, vector<char>& result) {
    ifstream ifs(filename, ios::binary|ios::ate);
    ifstream::pos_type pos = ifs.tellg();

    result.clear();//empty out
    result.resize(pos); //make sure we have enough space
    ifs.seekg(0, ios::beg);
    ifs.read(&result[0], pos);
}

char* createDigest(char* path) {
    stringstream ss;
    ss << "lzjd:" << path << ":";

    vector<char> all_bytes;
    readAllBytes(path, all_bytes);

    #ifdef DEBUGHASHES
    MurmurHash3 running_hash = MurmurHash3();
    int8_t some_byte = all_bytes[0];
    int32_t hash = running_hash.pushByte(some_byte);
    cout << "C++ MurmurHash3(" << (char)some_byte << ") = " << hash << endl;

    some_byte = all_bytes[1];
    hash = running_hash.pushByte(some_byte);
    cout << "C++ MurmurHash3(" << (char)some_byte << ") = " << hash << endl;
    #endif

    vector<int32_t> di = digest(digest_size, all_bytes);

    #ifdef DEBUGHASHES
    cout << "C++ MurmurHash ints:" << endl;
    for(unsigned int i = 0; i < di.size(); i++) {
        cout << di[i] << " ";
    }
    cout << endl;
    #endif

    typedef 
        bi::base64_from_binary<    // convert binary values to base64 characters
            bi::transform_width<   // retrieve 6 bit integers from a sequence of 32 bit ints
                vector<int32_t>::const_iterator,
                6,
                32
            >
        > 
        base64_text; 

    copy(
            base64_text(di.cbegin()),
            base64_text(di.cend()),
            ostream_iterator<char>(ss)
            );

    string lzjd_hash = ss.str();
    char* c_lzjd_hash = (char*) malloc(sizeof(char) * lzjd_hash.size());
    strcpy(c_lzjd_hash, lzjd_hash.c_str());
    return c_lzjd_hash;
}

char* createDigestFromBuffer(char *buff, int buffLen) {
    stringstream ss;
    ss << "lzjd:buffer:";

    vector<char> all_bytes;
    for(int i = 0; i < buffLen; i++) {
        all_bytes.push_back(buff[i]);
    }

    vector<int32_t> di = digest(digest_size, all_bytes);
    #ifdef DEBUGHASHES
    cout << "C++ MurmurHash ints:" << endl;
    for(unsigned int i = 0; i < di.size(); i++) {
        cout << di[i] << " ";
    }
    cout << endl;
    #endif

    typedef
        bi::base64_from_binary<    // convert binary values to base64 characters
            bi::transform_width<   // retrieve 6 bit integers from a sequence of 32 bit ints
                vector<int32_t>::const_iterator,
                6,
                32
            >
        >
        base64_text;

    copy(
            base64_text(di.cbegin()),
            base64_text(di.cend()),
            ostream_iterator<char>(ss)
            );

    string lzjd_hash = ss.str();
    char* c_lzjd_hash = (char*) malloc(sizeof(char) * lzjd_hash.size());
    strcpy(c_lzjd_hash, lzjd_hash.c_str());
    return c_lzjd_hash;
}

} // End Extern C
