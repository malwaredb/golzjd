#ifndef LZJD_HELPER_H
#define LZJD_HELPER_H

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

int32_t lzjd_similarity(char *hash1, char *hash2);
char* createDigest(char* path);
char* createDigestFromBuffer(char *buff, int buffLen);
    
#ifdef __cplusplus
}
#endif

#endif
