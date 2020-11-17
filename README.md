# workout-helper

Small tool to help with workout timing, and create a friendlier environment with a trainer voice.

## Back story

In spring in the the lockdown I worked out at home with plans that my kettlebel instructor sent me.

The workouts followed the similar patterns, with only a little variance.
Warmup at the start and stretch at the end were fix. In the middle there were sections with different approximate times. Once per week we had a tabata session with static exercises.

Now that the pandemic lockdown has arrived again, I come prepared for these workout sessions.

I'm happy with the result, because I do not need to check the time again, and a calm voice helps me through these workouts.

## Getting started

Currently I do not ship binaries, or docker image, so in order to run it you need to have go, I recommend the latest version.

The base workout consists of 10min warmup, 50 minutes of single task, 10 minutes of stretching.

```bash
go run .
```

The code can accept one string parameter. Tabata is marked with `t`, and general sections are marked with `s\d+`.
They have to be separated by a space character. Order is important.

e.g: `'t s10 s20'` means:

1. Warmup
1. A Tabata session
1. 10 minutes for one task
1. 20 minutes for another task
1. Stretching

```bash
go run . 't s10 s20'
```
