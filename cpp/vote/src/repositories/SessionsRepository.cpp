//
// Created by zede on 20.05.2023.
//

#include "SessionsRepository.h"


void SessionsRepository::addSession(const std::shared_ptr<Session>& session) {
  sessions_.try_emplace(session->id, session);
}

std::shared_ptr<Session> SessionsRepository::findSession(const std::string& id) const {
  auto result = sessions_.find(id);
  if (result == sessions_.end()) {
    return nullptr;
  }

  return result->second;
}