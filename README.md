<br>
<div align="center">
    <img src="./docs/logo-syntax-framework.png" />
    <p align="center">
        The ultimate web framework to do amazing things and be productive.
    </p>    
</div>

<br>

[//]: # ([![Build Status]&#40;https://github.com/syntax-framework/syntax/workflows/CI/badge.svg&#41;]&#40;https://github.com/syntax-framework/syntax/actions/workflows/ci.yml&#41; )
[//]: # ([![Documentation]&#40;https://img.shields.io/badge/documentation-gray&#41;]&#40;https://syntax-framework.com&#41;)

## WORK IN PROGRESS

- Declarative, html based
- Components, Directives
- Html Over the Wire (Phoenix LiveView like, less verbose)
- Http2
- No Javascript need (optional, but capable and powerful)
- Standalone (.exe) & Embedded (`go get -u syntax`)

```html
<component 
    name="clock" 
    param-title="string"
    js-param-label="@title"
>
  <h1>!{title} server side</h1>    
  <span>{{time}} - {{label}} client side</span>

  <style> span { color: red } </style>  

  <script>  
    const [time, setTime] = STX.state(new Date())

    const interval = setInterval(() => {
      setTime(new Date())  
    }, 1000)

    const clear = () => clearInterval(interval)  
  </script>  
</component>  

<clock title="My First Component"/>
```
