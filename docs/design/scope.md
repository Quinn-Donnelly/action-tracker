# Action Tracker

## Background

My thought here is to have a system that will get events about actions
and allow me to plot them out over time in various controlled fashion. The
motivation behind this is to have a system that would allow me to do the
following:

1. Have ownership of my data

2. Control the "met" criteria on expected days done[^1]

3. Customize the view it gives me for dopamine[^2]

## Tech

For the stack my current thought is to test out using [[TimeScaleDB]] for
the data storage. I'll use Golang as the back end serving the
functionality through API. For the view side of things I'm likely sticking
to the web for the first iteration of this but also focusing on getting
the data stores and served before locking in how to build the views up.

[^1]: I really want to have a system like [[Readwise]] for more of my
    actions that I'm tracking. Where it tracks days done and allows the
    recovery of your streak.

[^2]: I currently use [[Habits]] on android but it doesn't let me do quota
    based goals in a heatmap view. It always marks days as bad if they
    don't fuffil quota thus not supporting week based quota
