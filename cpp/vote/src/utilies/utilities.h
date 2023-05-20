//
// Created by zede on 20.05.2023.
//

#ifndef VOTE_SRC_UTILIES_UTILITIES_H_
#define VOTE_SRC_UTILIES_UTILITIES_H_

#include <string>
#include <random>
#include <sstream>
#include "restinio/all.hpp"

std::string generate_uuid();

template<typename T>
std::ostream &operator<<(std::ostream &to, const restinio::optional_t<T> &v) {
  if (v) to << *v;
  return to;
}

#endif //VOTE_SRC_UTILIES_UTILITIES_H_
