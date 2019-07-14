import {
  expect,
} from 'chai';
import {
  createChatroom,
  createUser,
  chatroomUsers,
  joinChatroom,
  deleteAllUsers,
} from './utils';

const newUser = (conn, name) => {
  return {
    conn,
    name,
    messages: [],
  };
};

/* eslint-disable func-names */
describe('User story', function () {
  const johnName = 'John';
  const masterName = 'Master';

  let chatroomID;
  let john;
  let master;

  const question1 = 'I have a question';
  const response1 = 'Huh?';
  const question2 = 'What\'s the best programming language in the world?';
  const response2 = 'PHP is the best language in the world';

  before(async function () {
    await deleteAllUsers();
    chatroomID = await createChatroom();

    await createUser(johnName);
    await createUser(masterName);

    const masterConn = joinChatroom(chatroomID, masterName);
    master = newUser(masterConn, masterName);
    const johnConn = joinChatroom(chatroomID, johnName);
    john = newUser(johnConn, johnName);
  });

  describe('John ask and Master response', function () {
    it('Messages are broadcasted and chatroom users list is correct', (done) => {
      try {
        john.conn.on('open', () => {
          john.conn.on('message', (data) => {
            const m = JSON.parse(data);
            john.messages.push(m.message);
          });

          setTimeout(() => {
            john.conn.send(question1);
            setTimeout(() => {
              john.conn.send(question2);
            }, 500);
          }, 0);
        });

        master.conn.on('open', () => {
          master.conn.on('message', (data) => {
            const m = JSON.parse(data);
            master.messages.push(m.message);
          });

          setTimeout(() => {
            john.conn.send(response1);
            setTimeout(() => {
              john.conn.send(response2);
            }, 500);
          }, 500);
        });

        setTimeout(async () => {
          try {
            const totalUsers = await chatroomUsers(chatroomID);
            expect(totalUsers).to.have.members([johnName, masterName]);

            expect(john.messages).to.have
              .members([question1, response1, question2, response2]);
            expect(master.messages).to.have
              .members([question1, response1, question2, response2]);

            john.conn.terminate();
            master.conn.terminate();
            done();
          } catch (e) {
            john.conn.terminate();
            master.conn.terminate();
            done(e);
          }
        }, 5000);
      } catch (e) {
        john.conn.terminate();
        master.conn.terminate();
        done(e);
      }
    });
  });

  describe('Rejoin chatroom to see history messages', function () {
    before(async function () {
      const johnConn = joinChatroom(chatroomID, johnName);
      john = newUser(johnConn, johnName);
    });

    it('Expect history messages exists', (done) => {
      try {
        john.conn.on('open', () => {
          john.conn.on('message', (data) => {
            const m = JSON.parse(data);
            john.messages.push(m.message);
          });
        });

        setTimeout(async () => {
          try {
            expect(john.messages).to.have
              .members([question1, response1, question2, response2]);
            john.conn.terminate();
            done();
          } catch (e) {
            john.conn.terminate();
            done(e);
          }
        }, 5000);
      } catch (e) {
        john.conn.terminate();
        done(e);
      }
    });
  });
});
/* eslint-enable func-names */
