const app = new Vue({
  el: '#elf-app',
  data: {
    // common
    documentScrollTop: 0,
    // post
    postForm: false,
    postPassword: '',
    postComments: [],
    postComment: {},
    postCommentParentLevel: null,
    postCaptchaUrl: null,
    postDirectories: []
  },
  computed: {
    pageKind () {
      return window.params.kind
    },
    pageData () {
      return window.params.data
    }
  },
  watch: {
    documentScrollTop (v) {
      let top = v + 5
      for (let i in this.postDirectories) {
        this.postDirectories[i].active = false
      }
      let target = -1
      for (let i in this.postDirectories) {
        if (top < this.postDirectories[i].top) {
          target = i - 1
          break
        }
      }
      if (target > -1) {
        this.postDirectories[target].active = true
      }
    }
  },
  methods: {
    // common
    datetimeFormat (date) {
      return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
    },
    gotop () {
      $('html,body').animate({scrollTop: 0}, 500)
    },
    // post
    postGetContent (event) {
      event.preventDefault()

      let uniqueId = this.pageData.uniqueId
      let password = this.postPassword
      let time = Math.ceil(new Date().getTime() / 60000)
      let ticket = md5(`${password}:${time}`)
      axios.get(`/content/${uniqueId}?ticket=${ticket}`).then(res => {
        let content = res.data.content
        this.postForm = false
        $('#content-body').html(content)
        this.$nextTick(() => {
          this.postContentInit()
        })
        // initCopy()
      }).catch(err => {
        alert('Incorrect password!')
      })
    },
    postGetComments () {
      let uniqueId = this.pageData.uniqueId
      axios.get(`/comment/${uniqueId}`).then(res => {
        this.postComments = res.data
        this.$nextTick(() => {
          if (location.hash) {
            let target = $(`a[name="${location.hash.substring(1)}"]`)
            if (target.length === 1) {
              $('html,body').scrollTop(target.offset().top)
            }
          }
        })
      })
    },
    postRefreshCaptcha () {
      this.postCaptchaUrl = `/captcha?_t=${new Date().getTime()}`
    },
    postSubmitComment (event) {
      event.preventDefault()

      let uniqueId = this.pageData.uniqueId
      let captcha = this.postComment.captcha
      axios.post(`/comment/${uniqueId}?captcha=${captcha}`, {
        parentLevel: this.postCommentParentLevel,
        nickname: this.postComment.nickname,
        email: this.postComment.email,
        content: this.postComment.content
      }).then(res => {
        $('#comment-modal').modal('hide')
        this.postGetComments()
      }, err => {
        alert('Error: ' + err)
      })
    },
    postGenerateDirectories () {
      let maxLevel = 6
      let space = '\u3000'
      let spaces = [space, space, space, space, space]
      this.postDirectories = $('#content-body').find('h1,h2,h3,h4,h5,h6').toArray().map(el => {
        let level = parseInt($(el).prop('tagName')[1])
        if (level < maxLevel) {
          maxLevel = level
        }
        let title = $(el).text()
        let top = $(el).offset().top
        return () => {
          return { top, active: false, text: spaces.slice(0, level - maxLevel).join('') + title }
        }
      }).map(_ => _())
    },
    postJumpDirectory (directory) {
      $('html,body').animate({scrollTop: directory.top}, 500)
    },
    postContentInit () {
      this.postGenerateDirectories()
      hljs.highlightAll()
    },
    // init
    initCommon () {},
    initPost () {
      this.postGetComments()
      this.postForm = this.pageData.private
      $('#comment-modal').on('show.bs.modal', () => {
        this.postRefreshCaptcha()
      })
      $('#comment-modal').on('hidden.bs.modal', () => {
        this.postCommentParentLevel = null
        this.postComment = {}
      })
      if (!this.pageData.private) {
        this.postContentInit()
      }
    }
  },
  mounted () {
    // before init

    // init
    this.initCommon()
    console.log(this.pageKind)
    switch (this.pageKind) {
    case 'post':
      this.initPost()
      break
    /* no default */
    }

    // after init
    window.addEventListener('scroll', () => {
      this.documentScrollTop = window.document.documentElement.scrollTop
    })
  }
})
