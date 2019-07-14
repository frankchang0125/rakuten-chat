import WebSocket from 'ws';
import superagent from 'superagent';
import logger from './logger';

const host = process.env.HOST;
const port = process.env.PORT;
const httpUrl = `http://${host}:${port}`;
const wsUrl = `ws://${host}:${port}`;

export const createChatroom = async () => {
  try {
    const url = `${httpUrl}/v1/chatrooms`;
    const resp = await superagent.post(url);
    return resp.body.id;
  } catch (e) {
    logger.error(`Cannot create chatroom: ${e}`);
    throw e;
  }
};

export const chatrooms = async () => {
  try {
    const url = `${httpUrl}/v1/chatrooms`;
    const resp = await superagent.get(url);
    return resp.body.ids;
  } catch (e) {
    logger.error(`Cannot get chatrooms list: ${e}`);
    throw e;
  }
};

export const chatroomUsers = async (chatroomID) => {
  try {
    const url = `${httpUrl}/v1/chatrooms/${chatroomID}/users`;
    const resp = await superagent.get(url);
    return resp.body.users;
  } catch (e) {
    logger.error(`Cannot get chatroom ${chatroomID} users list: ${e}`);
    throw e;
  }
};

export const createUser = async (userName) => {
  try {
    const url = `${httpUrl}/v1/users`;
    const body = {
      name: userName,
    };
    await superagent.post(url).send(body);
  } catch (e) {
    logger.error(`Cannot create user: ${e}`);
    throw e;
  }
};

export const joinChatroom = (chatroomID, userName) => {
  try {
    const url = `${wsUrl}/v1/chatrooms/${chatroomID}/join?user=${userName}`;
    const ws = new WebSocket(url);
    return ws;
  } catch (e) {
    logger.error(`Cannot join chatroom ${chatroomID}: ${e}`);
    throw e;
  }
};

export const deleteAllUsers = async () => {
  try {
    const url = `${httpUrl}/v1/users`;
    await superagent.delete(url);
  } catch (e) {
    logger.error(`Cannot delete all users: ${e}`);
    throw e;
  }
};
