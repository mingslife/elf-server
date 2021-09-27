'use strict';

$(() => {
  // 0. common
  let htmlEscape = text => text.replace(/[<>"&]/g, s => {
    switch (s) {
    case "<": return "&lt;"
    case ">": return "&gt;"
    case "&": return "&amp;"
    case "\"": return "&quot;"
    }
  })

  // 1. fetch content
  let fetchContent = () => {
    let uniqueId = window.uniqueId
    let password = $('#content-password').val()
    let time = Math.ceil(new Date().getTime() / 60000)
    let ticket = md5(`${password}:${time}`)
    axios.get(`/content/${uniqueId}?ticket=${ticket}`).then(res => {
      let content = res.data.content
      $('#content-form').remove()
      $('#content-body').html(content)
      initCopy()
    }).catch(err => {
      alert('Incorrect password!')
    })
  }
  $('#content-button').click(fetchContent)
  $('#content-password').keyup((e) => {
    if (e.keyCode === 13) {
      fetchContent()
    }
  })

  // 2. fetch comments
  let fetchComments = () => {
    let uniqueId = window.uniqueId
    axios.get(`/comment/${uniqueId}`).then(res => {
      let html = ''
      for (let i in res.data) {
        let comment = res.data[i]
        let commentedAtFormatted = dayjs(comment.commentedAt).format('YYYY-MM-DD HH:mm:ss')
        let nicknameEscaped = htmlEscape(comment.nickname)
        if (comment.username) {
          nicknameEscaped = `<a href="/user/${comment.username}">${nicknameEscaped}</a>`
        }
        let contentEscaped = comment.isBlocked ? '<span class="text-muted">BLOCKED</span>' : htmlEscape(comment.content).replace(/\t/g, '&nbsp;&nbsp;').replace(/\n/g, '<br />')
        let parentInfo = ''
        if (comment.parentLevel) {
          let parentNicknameEscaped = htmlEscape(comment.parentNickname)
          if (comment.parentUsername) {
            parentNicknameEscaped = `<a href="/user/${comment.parentUsername}">${parentNicknameEscaped}</a>`
          }
          parentInfo = `
<span>To</span>
<strong>${parentNicknameEscaped}</strong>
<span>(<a class="text-muted" href="#comment-${comment.parentLevel}">#${comment.parentLevel}</a>):</span>
          `
        }
        html +=
`
<div class="elf-comment">
  <div class="elf-comment__meta">
    <a class="text-muted" href="#comment-${comment.level}" name="comment-${comment.level}">#${comment.level}</a>
    <strong>${nicknameEscaped}</strong>
    <small class="text-muted">${commentedAtFormatted}</small>
    <div class="float-right elf-comment__actions">
      <a href="javascript:replyComment(${comment.level});"><i class="fa fa-reply"></i></a>
    </div>
  </div>
  <div class="elf-comment__content">
    <div>${parentInfo}</div>
    <div>${contentEscaped}</div>
  </div>
</div>
`
      }
      $('#comments-wrapper').html(html)

      setTimeout(() => {
        if (location.hash) {
          let target = $(`a[name="${location.hash.substring(1)}"]`)
          if (target.length === 1) {
            $('html,body').scrollTop(target.offset().top)
          }
        }
      }, 0)
    })
  }
  let refreshCommentCaptchaImage = () => $('#comment-captcha-image').attr('src', `/captcha?_t=${new Date().getTime()}`)
  $('#comment-modal').on('show.bs.modal', () => {
    refreshCommentCaptchaImage()

    $('#comment-modal-title').text(window.commentParentLevel ? 'Reply to #' + window.commentParentLevel : 'Comment')
  })
  $('#comment-modal').on('hidden.bs.modal', () => {
    window.commentParentLevel = null
    $('#comment-nickname').val('')
    $('#comment-email').val('')
    $('#comment-content').val('')
    $('#comment-captcha').val('')
  })
  $('#comment-captcha-image').click(refreshCommentCaptchaImage)
  window.commentSubmit = () => {
    let uniqueId = window.uniqueId
    let captcha = $('#comment-captcha').val()
    axios.post(`/comment/${uniqueId}?captcha=${captcha}`, {
      parentLevel: window['commentParentLevel'] || null,
      nickname: $('#comment-nickname').val(),
      email: $('#comment-email').val(),
      content: $('#comment-content').val()
    }).then(res => {
      $('#comment-modal').modal('hide')
      fetchComments()
    }, err => {
      alert('Error: ' + err)
    })
  }
  window.replyComment = parentLevel => {
    window.commentParentLevel = parentLevel
    $('#comment-modal').modal('show')
  }
  if (window['fetchComments']) {
    fetchComments()
  }

  // 3. highlight
  hljs.highlightAll()

  // 4. clipboard
  let initCopy = () => {
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
      text (trigger) {
        return $(trigger).parent().next().text()
      }
    }).on('success', e => {
      $(e.trigger).tooltip({
        title: 'OK',
        trigger: 'manual'
      }).tooltip('show')
      setTimeout(() => {
        $(e.trigger).tooltip('hide')
      }, 1500)
    })
  }
  initCopy()

  // directory
  let maxLevel = 6
  let space = '&nbsp;&nbsp;&nbsp;&nbsp;'
  let spaces = [space, space, space, space, space]
  let offsets = [{
    id: null,
    top: 0
  }]
  let menu = $('#content-body').find('h1,h2,h3,h4,h5,h6').toArray().map(element => {
    let level = parseInt($(element).prop('tagName')[1])
    if (level < maxLevel) {
      maxLevel = level
    }
    let title = $(element).text()
    let id = title.replace(/ /g, '-')
    $(element).attr('name', id)
    offsets.push({
      id: id,
      top: $(element).offset().top
    })
    return () => {
      return `<button class="dropdown-item" type="button" data-elf-menu="${id}">${spaces.slice(0, level - maxLevel).join('') + title}</button>`
    }
  }).map(_ => _()).join('')
  if (menu !== '') {
    $('#elf-corner__menu-group').show()
  }
  $('#elf-corner__menu').html(menu)
  $('#elf-corner__menu .dropdown-item').click(function() {
    let id = $(this).attr('data-elf-menu')
    let top = 0
    for (let i in offsets) {
      let offset = offsets[i]
      if (offset.id === id) {
        top = offset.top
        break
      }
    }
    $('html,body').animate({scrollTop: top}, 500)
  })

  $('#gotop-button').click(() => {
    $('html,body').animate({scrollTop: 0}, 500)
  })

  $(window).scroll(() => {
    let top = $(window).scrollTop() + 25
    let target = null
    for (let i in offsets) {
      let current = offsets[i]
      target = current.id
      if (i < offsets.length - 1) {
        let next = offsets[parseInt(i) + 1]
        if (top > current.top && top < next.top) {
          break
        }
      }
    }
    // console.log(target)
    $('#elf-corner__menu .dropdown-item').removeClass('active')
    $(`#elf-corner__menu .dropdown-item[data-elf-menu="${target}"]`).addClass('active')
  })
})
