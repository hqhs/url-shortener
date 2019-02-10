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
this is not the fastest available option. What we need is key/value store,
which is reliable out-of-the-box (hence not redis), easy to scale (simple
clustering), and easy to get up and running.

 Remember, what statistics above taken from startup which is **primary** product
is URL shortening & analytics service. Just copy-pasting their architecture
(and choosing hadoop) solutions doesn't fit because we shouldn't **expect**
such traffic. Our goal is to find balance between scale fees and time to get
service up and running.
 
 So, these are the options we have so far: AerospikeDB, Apache Cassandra, Riak,
DynamoDB (those are popular, there're other choices). I'm nearly not competent
enough for building infrastructure capable to withstand **real** high load
traffic (though some day I will), so let's just implement database connection
interface with opportunity to easy switch current choice in the future.
 
## URL shortener algorithm choice

Hash or not to hash?


**Variant 0: database sequence**:

Let short(url) be a function what returns different number for different urls,
which we base64 encode later. We want our service to start with N chars for one
url, so short(url) > 64^N + 1. We could use counter starting with 64^N + 1,
increment it after each url, convert it to base64 encoding.

Issue: It's easy to iterate over all urls. 

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

**Variant 1: Hashing**: 
hash(url, timestamp** and get N chars from result, make sure it's unique, get
other N chars if it's not.

Issue 0: It's slows down over time. If we used 1/x of all sequences, then each xth 
generation returns already used sequence and requires to query database again.

Issue 1: Possible race condition after large scaling.

**Variant 2: Zookeper**:
Pre-generate unique keys and store them as sequence. 
Issue: We need separate service to sync shortener nodes, and think about single 
point of failure.

So, I'm going with hashing. Let's calculate speed of service degradation with
N=6 on bit.ly load rates.

Given N=6 we have 6^64 - 5^64 ~ 10^50 possible keys. with 6 * 10^8
request/month, that's 7.2 * 10^9 per year and ~ 10^11 per 15 yers.

After 15 year we used 10^11 / 10^50 = 1/10^39 possible keys. Hence each 
10^39 request service would make additional database query. I could live with it :) 

## Thottling

## Caching

## Load balancing

## Actual deployment configuration

## Possible improvements

