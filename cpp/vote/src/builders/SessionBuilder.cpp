//
// Created by zede on 20.05.2023.
//

#include "SessionBuilder.h"

std::shared_ptr<Session> SessionBuilder::create(const std::string& user_agent, const std::string& id, const std::shared_ptr<User>& user) const {
  return std::make_shared<Session>(user_agent, id, user);
}