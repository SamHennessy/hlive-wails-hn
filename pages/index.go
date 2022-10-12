package pages

import (
	"context"
	"time"

	l "github.com/SamHennessy/hlive"
	"github.com/SamHennessy/hlive/hlivekit"
	"github.com/alexferrari88/gohn/pkg/gohn"
	"github.com/alexferrari88/gohn/pkg/processors"
	"github.com/andanhm/go-prettytime"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	pageSize               = 25
	topicShowStoryComments = "comments"
	topicLoadMoreStories   = "more"
)

func Index(ctxWails context.Context) *l.Page {
	// Hacker News Client
	hn := gohn.NewClient(nil)

	page := l.NewPage()

	page.DOM().HTML().Add(
		hlivekit.InstallPubSub(hlivekit.NewPubSub()), l.Class("h-full w-full"))
	page.DOM().Head().Add(
		l.T("link", l.Attrs{"rel": "stylesheet", "href": "/css/app.css"}),
		l.T("script", l.Attrs{"src": "/js/app.js", "defer": ""}))

	page.DOM().Body().Add(
		l.Class("h-full w-full bg-gradient-to-br from-slate-900 via-cyan-200 to-yellow-300"),
		NewCommentsView(hn),
		NewStoryList(hn),
	)

	return page
}

type StoryList struct {
	*hlivekit.ComponentPubSub

	storyIDs   []*int
	storyIndex int
}

func NewStoryList(hn *gohn.Client) *StoryList {
	c := &StoryList{
		ComponentPubSub: hlivekit.CPS("div", l.Class("h-full p-2")),
	}
	list := hlivekit.List("div")

	moreBtn := hlivekit.CPS("button", l.Class("cursor-pointer mt-2 text-xl"), "More")
	moreBtn.SetMountPubSub(func(ctx context.Context, pubSub *hlivekit.PubSub) {
		moreBtn.Add(l.On("click", func(ctx context.Context, e l.Event) {
			pubSub.Publish(topicLoadMoreStories, nil)
		}))
	})

	c.Add(
		l.T("div", l.Class("rounded bg-stone-100/75 px-2 h-full flex flex-col"),
			l.T("div", l.Class("text-4xl flex-none"), "Top Stories"),
			l.T("div", l.Class("grow overflow-auto"),
				list, moreBtn,
			),
		),
	)

	for i := 0; i < pageSize; i++ {
		list.Add(NewItemLoading())
	}

	c.SetMount(func(ctx context.Context) {
		// Get the top 500 stories' IDs
		go func() {
			c.storyIDs, _ = hn.Stories.GetTopIDs(ctx)
			runtime.WindowSetTitle(ctx, "Top Stories")
			list.RemoveAllItems()
			addStories(c, list, hn)

			l.RenderComponent(ctx, c)
		}()
	})

	c.SetMountPubSub(func(ctx context.Context, pubSub *hlivekit.PubSub) {
		pubSub.SubscribeFunc(func(message hlivekit.QueueMessage) {
			addStories(c, list, hn)
		}, topicLoadMoreStories)
	})

	return c
}

func addStories(c *StoryList, list *hlivekit.ComponentList, hn *gohn.Client) {
	offset := c.storyIndex
	for ; c.storyIndex < pageSize+offset && c.storyIndex < len(c.storyIDs); c.storyIndex++ {
		list.Add(NewStoryListView(hn, *c.storyIDs[c.storyIndex]))
	}
}

func deStr(v *string) string {
	if v == nil {
		return ""
	}

	return *v
}

func deInt(v *int) int {
	if v == nil {
		return 0
	}

	return *v
}

func NewStoryListView(hn *gohn.Client, id int) *StoryListView {
	c := NewItemLoading()

	c.SetMount(func(ctx context.Context) {
		go func() {
			story, _ := hn.Items.Get(ctx, id) // TODO: errors

			comments := hlivekit.CPS("div", l.Class("ml-1 cursor-pointer"), deInt(story.Descendants), " comments")
			comments.SetMountPubSub(func(ctx context.Context, pubSub *hlivekit.PubSub) {
				comments.Add(l.On("click", func(ctx context.Context, e l.Event) {
					pubSub.Publish(topicShowStoryComments, story)
				}))

				go l.RenderComponent(ctx, comments)
			})

			c.box.Set(l.T("div", l.Class("mt-2"),
				l.C("div", l.Class("text-lg cursor-pointer"), l.On("click", func(ctx context.Context, e l.Event) {
					runtime.BrowserOpenURL(ctx, deStr(story.URL))
				}), deStr(story.Title)),
				l.T("div", l.Class("flex text-sm text-stone-600"),
					l.T("div", deInt(story.Score), " points"),
					l.T("div", l.Class("ml-1"), "by ", deStr(story.By)),
					l.T("div", l.Class("ml-1"), formatTime(deInt(story.Time))),
					comments,
				),
			))

			l.RenderComponent(ctx, c)
		}()
	})

	return c
}

func NewItemLoading() *StoryListView {
	c := &StoryListView{
		ComponentMountable: l.CM("div"),

		box: l.Box(l.T("div", l.Class("mt-2"),
			l.T("div", l.Class("animate-pulse bg-stone-200 w-3/4 rounded h-5")),
			l.T("div", l.Class("animate-pulse bg-stone-200 w-2/4 rounded h-2 mt-2")),
		)),
	}

	c.Add(c.box)

	return c
}

type StoryListView struct {
	*l.ComponentMountable

	box *l.NodeBox[*l.Tag]
}

type CommentsView struct {
	*hlivekit.ComponentPubSub
}

func NewCommentsView(hn *gohn.Client) *CommentsView {
	c := &CommentsView{
		ComponentPubSub: hlivekit.CPS("div", l.Class("absolute top-0 left-0 w-full h-full p-2 hidden")),
	}

	back := l.C("div", l.Class("cursor-pointer flex-none p-2 text-lg"),
		"Back", l.On("click", func(ctx context.Context, e l.Event) {
			c.Add(l.Class("hidden"))
			runtime.WindowSetTitle(ctx, "Top Stories") // TODO: make dynamic
		}))
	storyBox := l.Box(l.T("div"))

	commentsLoading := l.T("div", l.Class("grow overflow-auto"))
	for i := 0; i < 10; i++ {
		commentsLoading.Add(l.T("div", l.Class("mt-2"),
			l.T("div", l.Class("animate-pulse bg-stone-200 w-1/4 rounded h-2")),
			l.T("div", l.Class("animate-pulse bg-stone-200 w-3/4 rounded h-5 mt-2")),
		))
	}

	commentsBox := l.Box(l.T("div"))
	c.Add(
		l.T("div", l.Class("flex flex-col h-full px-2 bg-stone-100 rounded"),
			back,
			storyBox,
			commentsBox,
		))

	c.SetMountPubSub(func(ctx context.Context, pubSub *hlivekit.PubSub) {
		pubSub.SubscribeFunc(func(message hlivekit.QueueMessage) {
			c.Add(l.ClassOff("hidden"))

			story, ok := message.Value.(*gohn.Item)
			if !ok {
				return
			}

			runtime.WindowSetTitle(ctx, deStr(story.Title))

			storyBox.Set(l.T("div", l.Class("flex-none p-2"),
				l.C("div", l.Class("text-lg cursor-pointer"), l.On("click", func(ctx context.Context, e l.Event) {
					runtime.BrowserOpenURL(ctx, deStr(story.URL))
				}), deStr(story.Title)),
				l.T("div", l.Class("flex text-sm"),
					l.T("div", deInt(story.Score), " points"),
					l.T("div", l.Class("ml-1"), "by ", deStr(story.By)),
					l.T("div", l.Class("ml-1"), deInt(story.Descendants), " comments"),
				),
			))

			commentsBox.Set(commentsLoading) // TODO: not always loading, why?!

			go func() {
				newBox := l.T("div", l.Class("grow overflow-auto"))
				commentsBox.Set(newBox)
				commentsMap, err := hn.Items.FetchAllDescendants(ctx, story, processors.UnescapeHTML())
				if err != nil {
					l.LoggerDev.Err(err).Msg("hn.Items.FetchAllDescendants")

					return
				}

				storyWithComments := gohn.Story{
					Parent:          story,
					CommentsByIdMap: commentsMap,
				}

				if storyWithComments.Parent.Kids == nil {
					return
				}

				kids := *storyWithComments.Parent.Kids

				for i := 0; i < len(kids); i++ {
					comment, exist := commentsMap[kids[i]]
					if exist {
						newBox.Add(commentRecurse(commentsMap, comment))
					} else {
						l.LoggerDev.Debug().Int("id", kids[i]).Int("story", *story.ID).Msg("missing comment on story")
					}
				}

				l.RenderComponent(ctx, c)
			}()
		}, topicShowStoryComments)
	})

	return c
}

func commentRecurse(commentsMap gohn.ItemsIndex, comment *gohn.Item) *l.Tag {
	b := l.T("div", l.Class("m-2 ml-2"),
		l.T("div", l.Class("m-1 border-l-2 border-stone-200"),
			l.T("div", l.Class("px-2 text-sm"), "by ", deStr(comment.By), " ", formatTime(deInt(comment.Time))),
			l.T("div", l.Class("px-2"), l.HTML(deStr(comment.Text))),
		))

	if comment.Kids != nil {
		kids := *comment.Kids
		for i := 0; i < len(kids); i++ {
			comment, exist := commentsMap[kids[i]]
			if exist {
				b.Add(commentRecurse(commentsMap, comment))
			} else {
				l.LoggerDev.Debug().Int("id", kids[i]).Int("comment", *comment.ID).Msg("missing kid comment")
			}
		}
	}

	return b
}

func formatTime(t int) string {
	if t == 0 {
		return ""
	}

	return prettytime.Format(time.Unix(int64(t), 0))
}
