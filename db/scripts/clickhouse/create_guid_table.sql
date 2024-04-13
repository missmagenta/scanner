CREATE TABLE IF NOT EXISTS guid (
                                         Mqtt String,
                                         Invent String,
                                         UnitGUID String,
                                         MsgID String,
                                         Text String,
                                         Context String,
                                         Class String,
                                         Level Int32,
                                         Area String,
                                         Addr String,
                                         Block String,
                                         Type String,
                                         Bit Int32,
                                         InvertBit Int32
) ENGINE = MergeTree()
ORDER BY (UnitGUID, MsgID);