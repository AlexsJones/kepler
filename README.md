```
_  _  ____  ____  __    ____  ____
( )/ )( ___)(  _ \(  )  ( ___)(  _ \
 )  (  )__)  )___/ )(__  )__)  )   /
(_)\_)(____)(__)  (____)(____)(_)\_)
```

Designed to do very few things.

- Allow you to batch execute commands in submodules
- Allow you to replace git package deps with `file` refs e.g. `git@github.com/tools/coolpackage.git` becomes `file:coolpackage`
  This is to help facilitate building from a meta-repository for multiple microservices in node that depend on each other.
