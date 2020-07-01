;(function () {
  const searchBtn = document.getElementById('search-btn')
  const userInput = document.getElementById('user-input')
  const content = document.getElementById('data-content')
  const url = `/FindUserInfo`

  function renderUser(data) {
    return `
    <ul>
        <li>勇者ID：${data.UserID}</li>
        <li>勇者名稱：${data.UserName}</li>
        <li>稱號：${data.Title}</li>
        <li>等級：${data.Level}</li>
        <li>種族：${data.Race}</li>
        <li>職業：${data.Occupation}</li>
        <li>巴幣：${data.Balance}</li>
        <li>GP：${data.GP}</li>
    </ul>
    `
  }

  function notFound() {
    return `<h2>查無此人</h2>`
  }

  function fetchUser(userID) {
    if (userID) {
      return new Promise((resolve, reject) => {
        fetch(url + '?ID=' + userID)
          .then(res => {
            console.log(res)
            return res.json()
          })
          .then(res => {
            console.log(res)
            resolve(res)
          })
          .catch(err => {
            reject(err)
          })
      })
    }
  }

  searchBtn.addEventListener('click', () => {
    userInput.classList.toggle('warning', userInput.value === '')
    if (userInput.value) {
      fetchUser(userInput.value).then(res => {
        console.log(res)
        if (res.data.UserID === '') {
          content.innerHTML = notFound()
        } else {
          content.innerHTML = renderUser(res.data)
        }
      })
    } else {
      console.log('no text')
    }
  })
})()
