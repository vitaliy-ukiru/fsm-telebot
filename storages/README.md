# storages

Directory with Storage implementations
Available now:

## Memory storage

Simple in-memory storage with synchronized data access.
Data is stored using maps.

## Redis storage

_Base
repository_: [github.com/nacknime-official/fsm-telebot-redis-storage](https://github.com/nacknime-official/fsm-telebot-redis-storage)

Storage using _redis_ as backend. For storing data using _encoding/gob_.

In this repository, you will find only the git submodule for the current repository.
This is done so that there are no problems with addiction, although it currently exists.
More details about this solution are written
in [#5](https://github.com/vitaliy-ukiru/fsm-telebot/pull/5#issuecomment-1666682226)

**If you have issues or want to contribute go to base repository.**

## File storage

Данное хранилище сохраняется в файлы и может
восстанавливать своё состояние из файлов.

For universality, saving data in the format is moved to the Provider interface.
This allows you not to think at the storage level about how the data will be stored.
The file/provider sub-package implements providers for such formats as:
- json (+ pretty version)
- gob
- base64 based at any provider

### json providers
There are two providers for JSON: _Json_ and _PrettyJson_.
What is their difference?

By the name, you can say that the second makes JSON readable, and in part this is true.

Both providers can return indented JSON. 
All kinds of `json.Encoder` and `json.Decoder` configurations via _JsonSettings_.

However, there is another difference as well.

To preserve type correctness in files, the data is stored as `map[string][]byte`.
The encoded value is stored as the value.

The encoding/json package turns `[]byte` into a base64 string.
This is exactly what "repairs" _PrettyJson_.

**But it's not free.**
The structure is copied to the new one to keep the data safe.

Also, _PrettyJson_ can have backward compatibility with _Json_. 
See PrettyJson docs for details.


### base64
If more example than real provider.

This provider works over other provider.
It encodes overlays base64 over the result from base provider.

Scheme with json provider:
```
Encoding: base64_encode(json_encode(data))
Decoding: json_decode(base64_decode(data))
```