Classic service example: URL shortener
===

## Features

- simple analytics
- ready-to-use kubernetes cluster for db & services.

## TODO: Code coverage

## Used sources

- []

  

## Database choice

  Let's say average long URL size is 2kb. Additionally I want to record
statistics, such as Unix timestamp and IP request address for each click. Keep
in mind that I may add more statistics in future. 10 bytes for timestamp and 16
bytes for IP (with IPv6 support) sounds reasonable. Some might argue what
statistics is not required, but I're living in word of data driven advertising
and loosing opportunity to collect data in my opinion is not realistic.

bit.ly statistics is ([link]( http://highscalability.com/blog/2014/7/14/bitly-lessons-learned-building-a-distributed-system-that-han.html )): 
  - 6 * 10^9 clicks/month
  - 6 * 10^8 shortens/month,
  - avarage click/link is 10, but maximum click per link could be 10^6 or even more.
 
 6 * 10^9 * 26 bytes + 6 * 10^8 * 2 kb = 1.35 terabytes/month
 
 I need to choose database wisely. Take into account what data has high
read/write ratio, practically never change (so eventual consistency means just
consistency), not structured (until I add users in the future), and record size
is small (around 2kb).
 
 I dont need rdbms -- data is too simple, not column store -- there's no need
for searching/sorting, not graph db -- there's no graph, not document store --
this is not the fastest available option. What I need is key/value store, which
is reliable out-of-the-box, easy to scale, and easy to get up and running.

 So, for bit.ly comparable traffic I can use one of those: AerospikeDB, Apache
Cassandra, Riak, DynamoDB (those are most popular, there're other choices). 

 Remember, what statistics above taken from startup which is **primary** product
is URL shortening and analytics service. Just copy-pasting their architecture
(and choosing hadoop) solutions doesn't fit because I can't **expect**
such traffic. And I'm not nearly competent enough for building infrastructure 
capable to withstand **real** high load traffic (though someday I will), so 
for purpose of making working demo I'm going with redis cluster and implement
database connection as golang interface, which allows us to quickly change
database then needed.

## URL shortener algorithm choice

Hash or not to hash?

**Variant 0: database sequence**:

Let short(url) be a function what returns different number for different urls,
which I base64 encode later. I want service to start with N chars for one
url, so short(url) > 64^N + 1. I could use counter starting with 64^N + 1,
increment it after each url, convert it to base64 encoding.

Issue: It's easy to iterate over all urls. 

**Variant 1: Hashing**: 
hash(url, timestamp** and get N chars from result, make sure it's unique, get
other N chars if it's not.

Issue 0: It's slows down over time. If I used 1/x of all sequences, then each xth 
generation returns already used sequence and requires to query database again.

Issue 1: Possible race condition after large scaling.

**Variant 2: Zookeper**:
Pre-generate unique keys and store them as sequence. 
Issue: I need separate service to sync shortener nodes, and think about single 
point of failure.

So, I'm going with hashing. Let's calculate speed of service degradation with
N=6 on bit.ly load rates.

Given N=6 I have 6^64 - 5^64 ~ 10^50 possible keys. with 6 * 10^8
request/month, that's 7.2 * 10^9 per year and ~ 10^11 per 15 years.

After 15 year I used 10^11 / 10^50 = 1/10^39 possible keys. Hence each 
10^39 request service would make additional database query. I could live with it :) 

## Possible improvements

### Thottling

### Caching

### Load balancing

- Telemetry
- Backpressure
- Health checking

### Other possible improvements

- Users and private statistics

## Actual deployment configuration with Kubernetes

## Contributing

