//
// Created by zede on 20.05.2023.
//

#include "user.h"

std::ostream &operator<<(std::ostream &to, const User &v) {
  to << fmt::format("{} ({}, {})", v.name, v.pwd, v.age);
  return to;
}