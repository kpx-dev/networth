import RxDB from 'rxdb';

RxDB.plugin(require('pouchdb-adapter-idb'));
const dbVersion = 0;

const networthSchema = {
  title: "networth schema",
  version: dbVersion,
  description: "hold networth info",
  type: "object",
  properties: {
    username: {
      type: "string",
      primary: true
    },
    networth: {
      type: "number",
    },
    updated_at: {
      type: "string",
    }
  }
};

const networthHistorySchema = {
  title: "networth history schema",
  version: dbVersion,
  description: "hold networth history info",
  type: "object",
  properties: {
    username: {
      type: "string",
      primary: true
    },
    networth: {
      type: "number",
    },
    updated_at: {
      type: "string",
    }
  }
};

let dbPromise = null;

const _create = async () => {
  const db = await RxDB.create({
    name: 'networth',
    adapter: 'idb',
  });

  await db.collection({name: 'networth', schema: networthSchema});
  await db.collection({name: 'networth_history', schema: networthHistorySchema});

  return db;
}

export const get = () => {
  if (!dbPromise)
      dbPromise = _create();
  return dbPromise;
}
