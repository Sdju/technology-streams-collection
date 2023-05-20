//
// Created by zede on 20.05.2023.
//

#ifndef VOTE_SRC_ENTITIES_COOKIE_H_
#define VOTE_SRC_ENTITIES_COOKIE_H_

#include <string>

struct Cookie {
 public:
  std::string name;
  std::string value;

  bool http_only = false;
  std::string path = "/";
  int max_age = 2592000;
};

#endif //VOTE_SRC_ENTITIES_COOKIE_H_
