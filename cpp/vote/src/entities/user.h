#ifndef VOTE_SRC_USER_H_
#define VOTE_SRC_USER_H_

#include <iostream>
#include "restinio/all.hpp"

struct User {
 public:
  std::string name;
  std::string pwd;
  int age;
};

std::ostream &operator<<(std::ostream &to, const User &v);

#endif //VOTE_SRC_USER_H_
