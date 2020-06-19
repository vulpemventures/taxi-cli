# taxi-cli
🖥 Command Line Interface to interact with 🚕 Liquid Taxi


## Install

**Download latest release**

* [Github Releases](https://github.com/vulpemventures/taxi-cli/releases)

Move and rename into a folder in your PATH (eg. `/usr/bin` or `/usr/local/bin`) and give permission eg. `chmod a+x ./taxi-cli`


## Commands
The following examples are for `regtest` local network. Just omit the `-regtest` flag to work with Liquid main network.

### Create

Create and print to console a base64 encoded PSET (Partial Signed Elements Transaction). It will fund a transaction with an input of `asset` greater than the given `amount`. It will add a change output to `from` address and an outout for `to` address. 

```sh
$ taxi-cli create -from ert1qnmlfhvhvy0zfnuptyl4rljys7kuxgdel94p0du -to ert1qeq52ql6cq6qs2j7xctylppzmy4aaekrxekqlt5 -amount 500000 -asset fd1a64da7a909b648319bc070e6f4d70ad1b515d70aa3183fcb137500e9670be -regtest
```
> NOTICE: You can specify a different Esplora REST endpoint to fetch your utxos. Eg. `-explorer https://myesplora.com`

### Topup

The transaction to be valid for broadcast needs L-BTC fees. Taxi allows you to do it without having them, you just provide a `pset` and pay in the given `asset`. Only a subset of Liquid assets are supported.

```sh
$ taxi-cli topup -pset cHNldP8BALgCAAAAAAETV88DBoA1hAe4uCSPNEpS5aPzDCy1qQPOl9V7n/bLlwEAAAAA/////wIBvnCWDlA3sfyDMapwXVEbrXBNbw4HvBmDZJuQetpkGv0BAAAAAAAHoSAAFgAUyCigf1gGgQVLxsLJ8IRbJXvc2GYBvnCWDlA3sfyDMapwXVEbrXBNbw4HvBmDZJuQetpkGv0BAADjX6kXgroAFgAUnv6bsuwjxJnwKyfqP8iQ9bhkNz8AAAAAAAEBQgG+cJYOUDex/IMxqnBdURutcE1vDge8GYNkm5B62mQa/QEAAONfqSfK6gAWABSe/puy7CPEmfArJ+o/yJD1uGQ3PwEDBAEAAAAAAAA= -asset fd1a64da7a909b648319bc070e6f4d70ad1b515d70aa3183fcb137500e9670be
```
> NOTICE: Outbound network connection is required to use the `topup` command

### Sign

You can sign a `pset`with a given `key`. Must be in hex format.

```sh
$ taxi-cli sign -pset cHNldP8BAP2RAQIAAAAAAhNXzwMGgDWEB7i4JI80SlLlo/MMLLWpA86X1Xuf9suXAQAAAAD/////r8/8HN9KPjIleplAQG5oOStCdJyTUci955FVxAojO9MBAAAAAP////8FAb5wlg5QN7H8gzGqcF1RG61wTW8OB7wZg2SbkHraZBr9AQAAAAAAB6EgABYAFMgooH9YBoEFS8bCyfCEWyV73NhmAb5wlg5QN7H8gzGqcF1RG61wTW8OB7wZg2SbkHraZBr9AQAA41+pF4K6ABYAFJ7+m7LsI8SZ8Csn6j/IkPW4ZDc/Ab5wlg5QN7H8gzGqcF1RG61wTW8OB7wZg2SbkHraZBr9AQAAAAAACKcQABYAFPojWKw/birmSNOMmgUCEaDPKhG0ASWyUQcOKcoZBDzzPM1zJOLdqwPsxK4LXnfE/A5c9slaAQAAAAAF9eDIABYAFPojWKw/birmSNOMmgUCEaDPKhG0ASWyUQcOKcoZBDzzPM1zJOLdqwPsxK4LXnfE/A5c9slaAQAAAAAAAAA4AAAAAAAAAAEBQgG+cJYOUDex/IMxqnBdURutcE1vDge8GYNkm5B62mQa/QEAAONfqSfK6gAWABSe/puy7CPEmfArJ+o/yJD1uGQ3PwEDBAEAAAAAAQFCASWyUQcOKcoZBDzzPM1zJOLdqwPsxK4LXnfE/A5c9slaAQAAAAAF9eEAABYAFPojWKw/birmSNOMmgUCEaDPKhG0IgIDb1ZG7WiLknk2naCkrXiVOufm0wBDbKijJkNg7+OCNuNIMEUCIQDdocVr3UYdCRFRjfNGgwasYk60byl8E3kek01eXcXlwQIgC7rwqfHib+E1ZUT5qXsxw2qJW+cSofqrbDuLF91iST4BAQMEAQAAAAAAAAAAAA== -key bfb96a215dfb07d1a193464174b9ea8e91f2a15bba79800dea838add330f6d86 -regtest
```

It will print both the finalized PSET and the extracted transaction as hex string ready to be broadcasted.









