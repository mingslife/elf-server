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
    postDirectories: [],
    postIsLogin: false,
    // reader
    readerEmail: '',
    readerValidateCode: '',
    readerValidateCodeTime: 0,
    readerValidateCodeTimer: null,
    readerCaptcha: '',
    readerCaptchaUrl: null,
    readerMode: 'info',
    readerRegisterData: {},
    readerInfoData: {},
    readerSendCoding: false
  },
  computed: {
    pageKind: () => window.params.kind,
    pageData: () => window.params.data
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
    toast (message) {
      setTimeout(() => {
        let toastId = 'elf-toast-' + new Date().getTime()
        $('#elf-toasts').prepend(`
        <div id="${toastId}" class="toast show rounded-0" role="alert" aria-live="assertive" aria-atomic="true" data-delay="5000">
          <div class="toast-body">${message}</div>
        </div>
        `)
        $(`#${toastId}`).toast('show').on('hidden.bs.toast', function() {
          $(this).remove()
        })
      }, 1)
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
      axios.post(`/comment/${uniqueId}`, {
        parentLevel: this.postCommentParentLevel,
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
      $('#content-body pre').each((_, element) => {
        $(element).before(`
        <div class="position-relative float-right elf-copy">
          <button type="button" class="btn bg-white btn-sm text-muted m-1 elf-copy__button">
            <i class="fa fa-copy"></i>
          </button>
        </div>
        `)
      })
      new ClipboardJS('.elf-copy__button', {
        text: (trigger) => $(trigger).parent().next().text()
      }).on('success', e => {
        $(e.trigger).tooltip({
          title: 'OK',
          trigger: 'manual'
        }).tooltip('show')
        setTimeout(() => {
          $(e.trigger).tooltip('hide')
        }, 1500)
      })
    },
    postToLogin () {
      location.href = '/reader?redirect=' + encodeURI(location.href)
    },
    initPost () {
      axios.get('/reader/info').then(res => {
        this.postIsLogin = true
      }).catch(err => {
        this.postIsLogin = false
      })
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
    },
    // reader
    readerRefreshCaptcha () {
      if (this.readerSendCoding) {
        return
      }
      this.readerCaptchaUrl = `/captcha?_t=${new Date().getTime()}`
    },
    readerShowCaptcha () {
      if (/^(\w)+(\.\w+)*@(\w)+((\.\w+)+)$/.test(this.readerEmail)) {
        $('#captcha-modal').modal('show')
      } else {
        this.toast(LOCALES.reader.wrongEmailAddress)
      }
    },
    readerGetInfo () {
      axios.get('/reader/info').then(res => {
        this.readerMode = 'info'
        this.readerEmail = res.data.email
        this.readerInfoData = {
          nickname: res.data.nickname,
          gender: res.data.gender,
          birthday: res.data.birthday ? dayjs(res.data.birthday).format("YYYY-MM-DD") : undefined,
          phone: res.data.phone
        }
      }).catch(err => {
        this.readerMode = 'login'
      })
    },
    readerUpdateInfo (event) {
      event.preventDefault()

      axios.post('/reader/info', {
        nickname: this.readerInfoData.nickname,
        gender: this.readerInfoData.gender,
        birthday: new Date(this.readerInfoData.birthday),
        phone: this.readerInfoData.phone
      }).then(res => {
        this.toast(LOCALES.reader.updateSuccessfully)
        this.readerGetInfo()
      }).catch(err => {
        this.toast(err.response.data.message)
      })
    },
    readerLogout () {
      axios.post('/reader/logout').finally(() => {
        this.readerInfoData = {}
        this.readerGetInfo()
      })
    },
    readerSendCode (event) {
      event.preventDefault()

      this.readerSendCoding = true
      axios.post('/reader/code', {
        email: this.readerEmail,
        captcha: this.readerCaptcha
      }).then(res => {
        this.readerValidateCodeTime = 60
        this.readerValidateCodeTimer = setInterval(() => {
          this.readerValidateCodeTime -= 1
          if (this.readerValidateCodeTime === 0) {
            clearInterval(this.readerValidateCodeTimer)
          }
        }, 1000)
        $('#captcha-modal').modal('hide')
      }).catch(err => {
        this.readerRefreshCaptcha()
        this.toast(err.response.data.message)
      }).finally(() => {
        this.readerSendCoding = false
      })
    },
    readerLogin (event) {
      event.preventDefault()

      axios.post('/reader/login', {
        email: this.readerEmail,
        validateCode: this.readerValidateCode
      }).then(res => {
        switch (res.data.result) {
        case 0:
          this.toast(LOCALES.reader.loginSuccessfully)
          this.readerEmail = ''
          this.readerValidateCode = ''
          this.readerMode = 'info'
          this.readerGetInfo()

          let params = new URLSearchParams(location.search)
          let redirect = params.get('redirect')
          if (redirect) {
            location.href = decodeURI(redirect)
          }

          break
        case 3:
          this.readerMode = 'register'
          break
        }
      }).catch(err => {
        this.toast(err.response.data.message)
      })
    },
    readerRegister (event) {
      event.preventDefault()

      axios.post('/reader/register', {
        email: this.readerEmail,
        validateCode: this.readerValidateCode,
        nickname: this.readerRegisterData.nickname,
        gender: this.readerRegisterData.gender,
        birthday: new Date(this.readerRegisterData.birthday),
        phone: this.readerRegisterData.phone
      }).then(res => {
        this.toast(LOCALES.reader.registerSuccessfully)
        this.readerEmail = ''
        this.readerValidateCode = ''
        this.readerRegisterData = {}
        this.readerMode = 'info'
        this.readerGetInfo()
      }).catch(err => {
        this.toast(err.response.data.message)
      })
    },
    initReader () {
      this.readerGetInfo()
      $('#captcha-modal').on('show.bs.modal', () => {
        this.readerRefreshCaptcha()
      })
      $('#captcha-modal').on('hidden.bs.modal', () => {
        this.readerCaptcha = ''
      })
    }
  },
  mounted () {
    // before init

    // init
    switch (this.pageKind) {
    case 'post':
      this.initPost()
      break
    case 'reader':
      this.initReader()
      break
    /* no default */
    }

    // after init
    window.addEventListener('scroll', () => {
      this.documentScrollTop = window.document.documentElement.scrollTop
    })
  }
})
