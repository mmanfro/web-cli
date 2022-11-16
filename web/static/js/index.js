document.getElementById("cmd").addEventListener(
    "keydown",
    (ev) => {
        if (ev.key === "Enter") {
            if (!ev.target.value) return

            const cmd = ev.target.value

            let history = document.getElementById("cmd-result")
            history.innerHTML = null
            document.getElementById("cmd-last").value = cmd
            ev.target.value = null

            fetch("/?cmd=" + encodeURIComponent(cmd))
                .then((response) => response.text())
                .then((text) => {
                    lines = text.split("\n")
                    lines.forEach((line) => {
                        history.appendChild(document.createTextNode(line))
                        history.appendChild(document.createElement("br"))
                    })
                })
            return
        }

        if (ev.key === "ArrowUp") {
            last_cmd = document.getElementById("cmd-last").value
            if (last_cmd)
                ev.target.value = document.getElementById("cmd-last").value
            return
        }
    },
    false
)

/**
 * _auth
 */

document.getElementById("auth").addEventListener("click", (ev) => {
    document.getElementById("modal-auth").style.display = "flex"
})

document.getElementById("modal-auth-close").addEventListener("click", (ev) => {
    document.getElementById("modal-auth").style.display = "none"
})

document.getElementById("login-set").addEventListener("click", (ev) => {
    let username = ev.target.querySelector("username")
    let password = ev.target.querySelector("password")
    if (ev.target.disabled) {
        ev.target.removeAttribute("disabled")
        username.removeAttribute("disabled")
        password.removeAttribute("disabled")
        username.value = null
        password.value = null
        document.cookie = "token=;"
        return
    }

    ev.target.setAttribute("disabled")
    username.setAttribute("disabled")
    password.setAttribute("disabled")
    username.value = null
    password.value = null
    document.cookie = "token=;"
    return
})
