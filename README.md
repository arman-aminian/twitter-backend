# A Twitter-like Website
![Booby Image](https://raw.githubusercontent.com/mehditeymorian/twitter-react/master/images/header_img.jpg)

This repository contains the backend code for the final project of Fall 2020 
Internet Engineering course. In this project, we were instructed to build a 
twitter-like website. We called our version "**Boobier**" after a
special seabird (more information [here](https://en.wikipedia.org/wiki/Booby)), 
and it's pretty much like Twitter. The frontend of this 
project is available [here](https://github.com/mehditeymorian/twitter-react).
The whole project is deployed on Heroku Cloud Application Platform.

## Project Explanation
We chose to write the backend in GoLang since it was the course syllabus,
and our database of choice was MongoDB because of its simplicity.

Here are the features that our website is capable of:
- Tweeting and deleting a tweet
- Making/Changing a profile (bio, profile picture, and header)
- Following/Unfollowing other users
- Notifications and Logs (history of a user's actions)
- Like, comment, and retweet
- Searching by username, text, and hashtag
- Viewing the timeline
- Seeing who liked and retweeted a tweet
- Hashtag trends
- Last but not least, user suggestion

### Database and Object Models
We currently have three databases: one for users, another for tweets, and
the last one for keeping the hashtags.

Until now, we have considered five objects to be modeled: users, tweets, owners
(explained later), hashtags, and events.

#### 1. User Model:
Each user `U` has the following fields:
- Name: name of `U`.
- Username: unique username.
- Email and password: email is considered unique and passwords are
hashed using `bcrypt` package.
- Bio, profile picture, and header: together, they form the user profile.
- Tweets: an array of `ObjectID`s where each id refers to a tweet.
- Followings/Followers: an array of `Owner`s which represents the profiles
`U` has followed or users that follow `U`.
- Notifications: an array of `Events` and it keeps track of three things:
    - One of `U`'s tweets were liked
    - One of `U`'s tweets were retweeted
    - Someone followed `U`
- Logs: Much like the Notifications, but it records `U`'s actions.

#### 2. Tweet Model:
Each tweet `T` is made of these fields:
- ID
- Text
- Media
- Date and Time of the tweet
- Owner: an `Owner` object which represents the tweet owner.
- Likes and Retweets: an array of `Owners`.
- Parent: generally empty but not when the tweet is a comment for another one.
- Comments: a list of `CommentTweet` objects.

The `CommentTweet` model contains all the `Tweet` fields except for the 
`Parent` and the `Comments`.

#### 3. Owner Model:
Each owner model is representing a user and only has the `Username`,
`Profile picture`, `Name`, `Bio`, and `IsFollowing` fields. The `IsFollowing`
field indicates whether the user requesting this object follows the target
user or not. For instance, all `U`'s followers has this field equal to
`true` for them.

#### 4. Hashtag Model:
Only records the name of a hashtag, the tweets which it belongs to, and
the number of times it was used in general.

#### 5. Event Model:
Three different actions are considered to be an `Event`: Like, Retweet, and
Following as explained before and each event has the following fields:
- Mode: whether it was a `Like`, a `Retweet`, or a `Follow` action.
- Source: the user causing the action.
- Target: the user to whom the action relates.
- Content: a short, simple description
- Timestamp
- Tweet: for `Like` and `Retweet` events shows the actual tweet.

The actual requests and corresponding responses can be seen in the code itself
and doesn't need much of an explanation.

##### Search by text
We tried to implement text-based search much like Twitter itself so our
search algorithm supports the following examples:
- "q1 q2 q3": all the tweets with this exact pattern in them.
- q1 q2 q3: all the tweets having at least one of the queries.
- "q1" q2 "q3": all the tweets that have q1 and q3. (q2 optional)

You can generate the docs of the backend to get a better sense of
requests and responses. To generate these docs (automatically with the
help of [swag](https://github.com/swaggo/swag)), first install the package
and then run ```swag init```, change the listening hostname to `localhost`
and then, the resulting documentation is available at 
[this link](http://localhost:8080/swagger/index.html).

## Contributions
There are many many bugs to be reported, suggestions to be told, ideas
to be shared, and other forms of feedback to be told. Any form of these
contributions would be a huge help to us improving this project. Thank
you in advance.

## Website and Team
You can test our website at [Boobier](https://booobier.herokuapp.com/).

Our Team:
- [Mohammad Mehdi Teymourian](https://github.com/mehditeymorian) (Frontend)
- [Arman Aminian](https://github.com/arman-aminian) (Backend)
- And me, [Mohammad Hosein Zarei](https://github.com/mhezarei) (Backend + 
a small part of frontend!)
