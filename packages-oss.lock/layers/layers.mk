# ***
# WARNING: Do not EDIT or MERGE this file, it is generated by packagespec.
# ***

LAYER_00-base-2017dceda969e0eab6fd00cab602c9a05538de9c_ID             := 00-base-2017dceda969e0eab6fd00cab602c9a05538de9c
LAYER_00-base-2017dceda969e0eab6fd00cab602c9a05538de9c_TYPE           := base
LAYER_00-base-2017dceda969e0eab6fd00cab602c9a05538de9c_BASE_LAYER     := 
LAYER_00-base-2017dceda969e0eab6fd00cab602c9a05538de9c_SOURCE_INCLUDE := 
LAYER_00-base-2017dceda969e0eab6fd00cab602c9a05538de9c_SOURCE_EXCLUDE := 
LAYER_00-base-2017dceda969e0eab6fd00cab602c9a05538de9c_CACHE_KEY_FILE := .buildcache/cache-keys/base-2017dceda969e0eab6fd00cab602c9a05538de9c
LAYER_00-base-2017dceda969e0eab6fd00cab602c9a05538de9c_ARCHIVE_FILE   := .buildcache/archives/00-base-2017dceda969e0eab6fd00cab602c9a05538de9c.tar.gz
$(eval $(call LAYER,$(LAYER_00-base-2017dceda969e0eab6fd00cab602c9a05538de9c_ID),$(LAYER_00-base-2017dceda969e0eab6fd00cab602c9a05538de9c_TYPE),$(LAYER_00-base-2017dceda969e0eab6fd00cab602c9a05538de9c_BASE_LAYER),$(LAYER_00-base-2017dceda969e0eab6fd00cab602c9a05538de9c_SOURCE_INCLUDE),$(LAYER_00-base-2017dceda969e0eab6fd00cab602c9a05538de9c_SOURCE_EXCLUDE),$(LAYER_00-base-2017dceda969e0eab6fd00cab602c9a05538de9c_CACHE_KEY_FILE),$(LAYER_00-base-2017dceda969e0eab6fd00cab602c9a05538de9c_ARCHIVE_FILE)))

LAYER_01-ui-f43e20c1d27527cddb9fa2f5a08420e5058396f1_ID             := 01-ui-f43e20c1d27527cddb9fa2f5a08420e5058396f1
LAYER_01-ui-f43e20c1d27527cddb9fa2f5a08420e5058396f1_TYPE           := ui
LAYER_01-ui-f43e20c1d27527cddb9fa2f5a08420e5058396f1_BASE_LAYER     := 00-base-2017dceda969e0eab6fd00cab602c9a05538de9c
LAYER_01-ui-f43e20c1d27527cddb9fa2f5a08420e5058396f1_SOURCE_INCLUDE := internal/ui/VERSION
LAYER_01-ui-f43e20c1d27527cddb9fa2f5a08420e5058396f1_SOURCE_EXCLUDE := 
LAYER_01-ui-f43e20c1d27527cddb9fa2f5a08420e5058396f1_CACHE_KEY_FILE := .buildcache/cache-keys/ui-f43e20c1d27527cddb9fa2f5a08420e5058396f1
LAYER_01-ui-f43e20c1d27527cddb9fa2f5a08420e5058396f1_ARCHIVE_FILE   := .buildcache/archives/01-ui-f43e20c1d27527cddb9fa2f5a08420e5058396f1.tar.gz
$(eval $(call LAYER,$(LAYER_01-ui-f43e20c1d27527cddb9fa2f5a08420e5058396f1_ID),$(LAYER_01-ui-f43e20c1d27527cddb9fa2f5a08420e5058396f1_TYPE),$(LAYER_01-ui-f43e20c1d27527cddb9fa2f5a08420e5058396f1_BASE_LAYER),$(LAYER_01-ui-f43e20c1d27527cddb9fa2f5a08420e5058396f1_SOURCE_INCLUDE),$(LAYER_01-ui-f43e20c1d27527cddb9fa2f5a08420e5058396f1_SOURCE_EXCLUDE),$(LAYER_01-ui-f43e20c1d27527cddb9fa2f5a08420e5058396f1_CACHE_KEY_FILE),$(LAYER_01-ui-f43e20c1d27527cddb9fa2f5a08420e5058396f1_ARCHIVE_FILE)))

LAYER_02-go-modules-0542ebe70aa49ea525c7cd85313d5ba08f90441e_ID             := 02-go-modules-0542ebe70aa49ea525c7cd85313d5ba08f90441e
LAYER_02-go-modules-0542ebe70aa49ea525c7cd85313d5ba08f90441e_TYPE           := go-modules
LAYER_02-go-modules-0542ebe70aa49ea525c7cd85313d5ba08f90441e_BASE_LAYER     := 01-ui-f43e20c1d27527cddb9fa2f5a08420e5058396f1
LAYER_02-go-modules-0542ebe70aa49ea525c7cd85313d5ba08f90441e_SOURCE_INCLUDE := go.mod go.sum */go.mod */go.sum
LAYER_02-go-modules-0542ebe70aa49ea525c7cd85313d5ba08f90441e_SOURCE_EXCLUDE := 
LAYER_02-go-modules-0542ebe70aa49ea525c7cd85313d5ba08f90441e_CACHE_KEY_FILE := .buildcache/cache-keys/go-modules-0542ebe70aa49ea525c7cd85313d5ba08f90441e
LAYER_02-go-modules-0542ebe70aa49ea525c7cd85313d5ba08f90441e_ARCHIVE_FILE   := .buildcache/archives/02-go-modules-0542ebe70aa49ea525c7cd85313d5ba08f90441e.tar.gz
$(eval $(call LAYER,$(LAYER_02-go-modules-0542ebe70aa49ea525c7cd85313d5ba08f90441e_ID),$(LAYER_02-go-modules-0542ebe70aa49ea525c7cd85313d5ba08f90441e_TYPE),$(LAYER_02-go-modules-0542ebe70aa49ea525c7cd85313d5ba08f90441e_BASE_LAYER),$(LAYER_02-go-modules-0542ebe70aa49ea525c7cd85313d5ba08f90441e_SOURCE_INCLUDE),$(LAYER_02-go-modules-0542ebe70aa49ea525c7cd85313d5ba08f90441e_SOURCE_EXCLUDE),$(LAYER_02-go-modules-0542ebe70aa49ea525c7cd85313d5ba08f90441e_CACHE_KEY_FILE),$(LAYER_02-go-modules-0542ebe70aa49ea525c7cd85313d5ba08f90441e_ARCHIVE_FILE)))

LAYER_03-copy-source-e59a64ba925a852a09535a8d252147336d635eff_ID             := 03-copy-source-e59a64ba925a852a09535a8d252147336d635eff
LAYER_03-copy-source-e59a64ba925a852a09535a8d252147336d635eff_TYPE           := copy-source
LAYER_03-copy-source-e59a64ba925a852a09535a8d252147336d635eff_BASE_LAYER     := 02-go-modules-0542ebe70aa49ea525c7cd85313d5ba08f90441e
LAYER_03-copy-source-e59a64ba925a852a09535a8d252147336d635eff_SOURCE_INCLUDE := *.go *.up.sql
LAYER_03-copy-source-e59a64ba925a852a09535a8d252147336d635eff_SOURCE_EXCLUDE := 
LAYER_03-copy-source-e59a64ba925a852a09535a8d252147336d635eff_CACHE_KEY_FILE := .buildcache/cache-keys/copy-source-e59a64ba925a852a09535a8d252147336d635eff
LAYER_03-copy-source-e59a64ba925a852a09535a8d252147336d635eff_ARCHIVE_FILE   := .buildcache/archives/03-copy-source-e59a64ba925a852a09535a8d252147336d635eff.tar.gz
$(eval $(call LAYER,$(LAYER_03-copy-source-e59a64ba925a852a09535a8d252147336d635eff_ID),$(LAYER_03-copy-source-e59a64ba925a852a09535a8d252147336d635eff_TYPE),$(LAYER_03-copy-source-e59a64ba925a852a09535a8d252147336d635eff_BASE_LAYER),$(LAYER_03-copy-source-e59a64ba925a852a09535a8d252147336d635eff_SOURCE_INCLUDE),$(LAYER_03-copy-source-e59a64ba925a852a09535a8d252147336d635eff_SOURCE_EXCLUDE),$(LAYER_03-copy-source-e59a64ba925a852a09535a8d252147336d635eff_CACHE_KEY_FILE),$(LAYER_03-copy-source-e59a64ba925a852a09535a8d252147336d635eff_ARCHIVE_FILE)))
