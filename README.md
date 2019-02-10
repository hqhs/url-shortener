Classic service example: URL shortener
===

## Features

- simple analytics

## TODO: Code coverage

## Database choice

Let's say average long URL size is 2kb. Additionally we want to record
statistics, such as Unix timestamp and IP request address for each click. Keep
in mind that we may add more statistics in future. 10 bytes for timestamp and 16
bytes for IP (with IPv6 support) sounds reasonable. Some might argue what
statistics is not required, but we're living in word of data driven advertising
and loosing opportunity to collect data in my opinion is not realistic.

bit.ly statistics is ([link]( http://highscalability.com/blog/2014/7/14/bitly-lessons-learned-building-a-distributed-system-that-han.html )): 
  - 6 * 10^9 clicks/month
  - 6 * 10^8 shortens/month,
  - avarage click/link is 10, but max(click per link) could be 10^6 or even more.
 
 6 * 10^9 * 26 bytes + 6 * 10^8 * 2 kb = 1.35 terabytes/month
 
 Obviously we need to choose database wisely. Out data structure is fixed,
 with small quantities, we don't need atomicity.
 
 We dont need rdbms -- data is too simple, not column store -- there's no need
 for searching/sorting, not graph db -- there's no graph, not document store -- 
 this is not the fastest option. What we need is key/value store, which is
 reliable out-of-the-box (hence not redis cluster), easy to scale (simple
 clustering), and easy to get up and running (**not** a hadoop cluster).

 Remember, what statistics above taken from startup which is **primary** product
 is URL shortening & analytics service. Just copy-pasting their architecture
 (and choosing hadoop) solutions doesn't fit because we shouldn't **expect**
 such traffic. Our goal is to find balance between scale fees and time to get
 service up and running.
 
## URL shortener algorithm choice

Hash or not to hash?

Let short(url) 
be a function what returns number, which we base64 encode later.
We want our service to start with 6 chars for one url, so should always be
short(url) > 64^5 + 1.

Variant 0: use counter starting with 64^5 + 1, increment it after each url, 
convert it with base64 encoding

``` python
def int2base64(x):
    'convert an integer to its url safe string representation in a given base'
    alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
    rets=''
    while x>0:
        x,idx = divmod(x,64)
        rets = alphabet[idx] + rets
    return rets
```

Variant 1: use  n


## statistics: codebase size & time spent

## possible improvements

